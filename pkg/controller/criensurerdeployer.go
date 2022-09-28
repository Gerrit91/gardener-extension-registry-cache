package controller

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/gardener/gardener/pkg/utils/imagevector"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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

	ReferencedServices *corev1.ServiceList
}

const (
	criEnsurerName = "cri-config-ensurer"
	scriptOfDeath  = `
#!/usr/bin/env bash

set -euo pipefail

CONTAINERD_IMPORTS_DIR="/etc/containerd/conf.d"
CONFIG_INPUT_FILE=$1

if ! grep -F '/etc/containerd/conf.d/*.toml' /host/etc/containerd/config.toml; then
	# https://github.com/gardener/gardener/blob/v1.51.0/docs/usage/custom-containerd-config.md
	echo "ERROR: Only works on workers created with Gardener >v1.51, exiting."
	exit 1
fi

if [ ! -e "$CONFIG_INPUT_FILE" ]; then
	echo "ERROR: Config input file $CONFIG_INPUT_FILE could not be found, exiting."
	exit 1
fi

mkdir -p "/host$CONTAINERD_IMPORTS_DIR"

while true; do
	input_file=$(cat "$CONFIG_INPUT_FILE")
	existing_file=$(cat "/host$CONTAINERD_IMPORTS_DIR/$(basename "$CONFIG_INPUT_FILE")")
	if [[ "$input_file" != "$existing_file" ]]; then
		echo "restarting containerd"
		cp -f "$CONFIG_INPUT_FILE" "/host$CONTAINERD_IMPORTS_DIR/"
		chroot /host systemctl restart containerd.service
		echo "applied registry mirrors, sleeping for a minute"
	else
		echo "no changes required, sleeping for a minute"
	fi
	sleep 60
done
`
)

func (c *criEnsurer) Ensure() ([]client.Object, error) {
	if c.Labels == nil {
		c.Labels = map[string]string{
			"app": c.Name,
		}
	}

	toml, err := c.configToml()
	if err != nil {
		return nil, fmt.Errorf("unable to template toml: %w", err)
	}

	var (
		configMap = &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      c.Name,
				Namespace: c.Namespace,
				Labels:    c.Labels,
			},
			Data: map[string]string{
				"reconcile.sh":                     scriptOfDeath,
				"zz-extension-registry-cache.toml": toml,
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
									"bash", "-c", "/work/reconcile.sh /work/zz-extension-registry-cache.toml",
								},
								ImagePullPolicy: v1.PullIfNotPresent,
								VolumeMounts: []v1.VolumeMount{
									{
										Name:      "work",
										ReadOnly:  true,
										MountPath: "/work",
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
								Name: "work",
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
	}, nil
}

func (c *criEnsurer) configToml() (string, error) {
	type criMirror struct {
		Host     string
		Endpoint string
	}

	var mirrors []*criMirror
	for i := range c.ReferencedServices.Items {
		svc := c.ReferencedServices.Items[i]
		mirrors = append(mirrors, &criMirror{
			Host:     svc.Labels[registryCacheServiceUpstreamLabel],
			Endpoint: fmt.Sprintf("http://%s:%d", svc.Spec.ClusterIP, svc.Spec.Ports[0].Port),
		})
	}

	text := `# governed by gardener-extension-registry-cache, do not edit
{{ range $mirror := . -}}
[plugins."io.containerd.grpc.v1.cri".registry.mirrors."{{ $mirror.Host }}"]
  endpoint = ["{{ $mirror.Endpoint }}"]
{{ end }}`

	tpl, err := template.New("").Parse(text)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, mirrors); err != nil {
		return "", err
	}

	return buf.String(), nil
}
