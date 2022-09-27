<p>Packages:</p>
<ul>
<li>
<a href="#service.registry.extensions.gardener.cloud%2fv1alpha1">service.registry.extensions.gardener.cloud/v1alpha1</a>
</li>
</ul>
<h2 id="service.registry.extensions.gardener.cloud/v1alpha1">service.registry.extensions.gardener.cloud/v1alpha1</h2>
<p>
<p>Package v1alpha1 contains the registry service extension.</p>
</p>
Resource Types:
<ul><li>
<a href="#service.registry.extensions.gardener.cloud/v1alpha1.RegistryConfig">RegistryConfig</a>
</li></ul>
<h3 id="service.registry.extensions.gardener.cloud/v1alpha1.RegistryConfig">RegistryConfig
</h3>
<p>
<p>RegistryConfig configuration resource</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
service.registry.extensions.gardener.cloud/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>RegistryConfig</code></td>
</tr>
<tr>
<td>
<code>registries</code></br>
<em>
<a href="#service.registry.extensions.gardener.cloud/v1alpha1.RegistryMirror">
[]RegistryMirror
</a>
</em>
</td>
<td>
<p>Mirrors is a slice of registry mirrors to deploy</p>
</td>
</tr>
</tbody>
</table>
<h3 id="service.registry.extensions.gardener.cloud/v1alpha1.RegistryMirror">RegistryMirror
</h3>
<p>
(<em>Appears on:</em>
<a href="#service.registry.extensions.gardener.cloud/v1alpha1.RegistryConfig">RegistryConfig</a>)
</p>
<p>
<p>RegistryMirror defines a registry mirror to deploy</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>remoteURL</code></br>
<em>
string
</em>
</td>
<td>
<p>RemoteURL is the remote URL of registry to mirror</p>
</td>
</tr>
<tr>
<td>
<code>port</code></br>
<em>
int
</em>
</td>
<td>
<p>Port is the port on which the registry mirror is going to serve</p>
</td>
</tr>
<tr>
<td>
<code>cacheSize</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>CacheSize is the size of the registry cache</p>
</td>
</tr>
<tr>
<td>
<code>cacheGarbageCollectionEnabled</code></br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>CacheGarbageCollectionEnabled enables/disables cache garbage collection</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<p><em>
Generated with <a href="https://github.com/ahmetb/gen-crd-api-reference-docs">gen-crd-api-reference-docs</a>
</em></p>
