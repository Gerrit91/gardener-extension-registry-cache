package controller

import (
	"fmt"
	"strings"

	"github.com/gardener/gardener/pkg/utils/imagevector"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type criEnsurer struct {
	Name      string
	Namespace string
	Labels    map[string]string

	CRIEnsurerImage *imagevector.Image

	RegistryMirrors map[string]string
}

const (
	criEnsurerName = "cri-config-ensurer"
	scriptOfDeath  = `
#!/bin/sh

set -euo pipefail

while true; do
	echo "applying registry mirrors"

	changed=false

	for mirror in $@; do
	registry=$(echo $mirror | cut -f1 -d'@')
	endpoint=$(echo $mirror | cut -f2 -d'@')

	line1="[plugins.\"io.containerd.grpc.v1.cri\".registry.mirrors.\"${registry}\"]"
	line2="  endpoint = [\"${endpoint}\"]"

	cat /host/etc/containerd/config.toml | grep -F "${line1}" && continue

	changed=true

	printf "$line1\n" >> /host/etc/containerd/config.toml
	printf "$line2\n" >> /host/etc/containerd/config.toml
	done

	if [ "$changed" = true ]; then
	echo "restarting containerd"
	chroot /host systemctl restart containerd

	echo "applied registry mirrors, sleeping for a minute"
	else
	echo "no changes required, sleeping for a minute"
	fi

	sleep 60
done
`
)

func (c *criEnsurer) Ensure() []client.Object {
	if c.Labels == nil {
		c.Labels = map[string]string{
			"app": c.Name,
		}
	}

	var registryMirrors []string
	for host, address := range c.RegistryMirrors {
		registryMirrors = append(registryMirrors, fmt.Sprintf(`'%s@%s'`, host, address))
	}

	var (
		configMap = &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      c.Name,
				Namespace: c.Namespace,
				Labels:    c.Labels,
			},
			Data: map[string]string{
				"reconcile.sh": scriptOfDeath,
			},
		}

		daemonSet = &appsv1.DaemonSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      c.Name,
				Namespace: registryCacheNamespaceName,
				Labels:    c.Labels,
			},
			Spec: appsv1.DaemonSetSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: c.Labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: c.Labels,
					},
					Spec: v1.PodSpec{
						HostPID: true,
						Containers: []v1.Container{
							{
								Name:  criEnsurerName,
								Image: c.CRIEnsurerImage.Repository,
								SecurityContext: &v1.SecurityContext{
									Privileged: pointer.Bool(true),
								},
								Command: []string{
									"sh", "-c", fmt.Sprintf("/scripts/reconcile.sh %s", strings.Join(registryMirrors, " ")),
								},
								ImagePullPolicy: v1.PullIfNotPresent,
								VolumeMounts: []v1.VolumeMount{
									{
										Name:      "script",
										ReadOnly:  true,
										MountPath: "/scripts",
									},
									{
										Name:      "host",
										MountPath: "/host",
									},
								},
							},
						},
						Volumes: []v1.Volume{
							{
								Name: "script",
								VolumeSource: v1.VolumeSource{
									ConfigMap: &v1.ConfigMapVolumeSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: c.Name,
										},
										DefaultMode: pointer.Int32(0744),
									},
								},
							},
							{
								Name: "host",
								VolumeSource: v1.VolumeSource{
									HostPath: &v1.HostPathVolumeSource{
										Path: "/",
									},
								},
							},
						},
					},
				},
			},
		}
	)

	return []client.Object{
		configMap,
		daemonSet,
	}
}
