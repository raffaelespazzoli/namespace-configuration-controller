# Namespace Configuration Controller

The namespace configuration controller helps keeping a namespace's configuration aligned with one of more policies specified as a CRD.
Currently the following objects are part of the namespace configuration:

- ConfigMaps         
- PodPresets         
- Quotas              
- LimitRanges        
- RoleBindings        
- ClusterRoleBindings 
- ServiceAccounts

Dev teams may of may not be granted permissions to create these objects. In case they haven't the namespace configuration controller can be a way to create namespace configuration policy and govern the way namespace are configured.

A NamespaceConfig CRD looks as follows:

```
apiVersion: namespaceconfig.raffaelespazzoli.systems/v1alpha1
kind: NamespaceConfig
metadata:
  name: example-namespaceconfig
spec:
  selector:
    matchLabels:
      namespaceconfig: "true"
    matchExpressions:
     - {key: namespaceconfig, operator: In, values: ["true"]}  
  networkpolicies: []
  configmaps: []         
  podpresets: []         
  quotas: []              
  limitranges: []        
  rolebindings: []        
  clusterrolebindings: [] 
  serviceaccounts: []
```

The selector will select the namespaces to which this configuration should be applied.
In this example all the managed ojects types have a empty array.
You can add your API object instance there. The namespace field should not be specified and if it exists it will be overwrittent with the namespace name of the namespace to which the configuration is being applied.

## Example of configurations

Here is a list of use cases in which the Namespace Configuration Controller can be useful

### T-Shirt sized quotas

during the provisionin gof ne projects to dev teams some organizations start with T-shirt sized quotas. Here is an example of how this can be done with the Namespace Configuration Controller

```
apiVersion: namespaceconfig.raffaelespazzoli.systems/v1alpha1
kind: NamespaceConfig
metadata:
  name: small-size
spec:
  selector:
    matchLabels:
      size: small  
  quotas:
  - apiVersion: v1
    kind: ResourceQuota
    metadata:
    name: small-size
    spec:
    hard:
      request.cpu: "4" 
      request.memory: "2Gi" 
      limits.ephemeral-storage: "4Gi"
---
apiVersion: namespaceconfig.raffaelespazzoli.systems/v1alpha1
kind: NamespaceConfig
metadata:
  name: large-size
spec:
  selector:
    matchLabels:
      size: large  
  quotas:
  - apiVersion: v1
    kind: ResourceQuota
    metadata:
    name: large-size
    spec:
    hard:
      request.cpu: "8" 
      request.memory: "4Gi" 
      limits.ephemeral-storage: "8Gi"      
```

We can test the above configuration as follows:
```
oc apply -f examples
oc new-project large-project
oc label namespace large-project size=large
oc new-project small-project
oc label namespace small-project size=small
```