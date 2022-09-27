<p>Packages:</p>
<ul>
<li>
<a href="#registry.extensions.gardener.cloud%2fv1alpha1">registry.extensions.gardener.cloud/v1alpha1</a>
</li>
</ul>
<h2 id="registry.extensions.gardener.cloud/v1alpha1">registry.extensions.gardener.cloud/v1alpha1</h2>
<p>
<p>Package v1alpha1 contains the registry service extension.</p>
</p>
Resource Types:
<ul></ul>
<h3 id="registry.extensions.gardener.cloud/v1alpha1.RegistryConfig">RegistryConfig
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
<code>mirrors</code></br>
<em>
<a href="#registry.extensions.gardener.cloud/v1alpha1.RegistryMirror">
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
<h3 id="registry.extensions.gardener.cloud/v1alpha1.RegistryMirror">RegistryMirror
</h3>
<p>
(<em>Appears on:</em>
<a href="#registry.extensions.gardener.cloud/v1alpha1.RegistryConfig">RegistryConfig</a>)
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
<code>upstreamURL</code></br>
<em>
string
</em>
</td>
<td>
<p>UpstreamURL is the remote URL of registry to mirror</p>
</td>
</tr>
<tr>
<td>
<code>port</code></br>
<em>
int32
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
