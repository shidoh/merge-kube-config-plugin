### Installing the plugin:

1. Save this code in a file named `kubectl-mergekubeconfig.go`.
2. Compile the program with the command:

```sh
go build -o kubectl-mergekubeconfig kubectl-mergekubeconfig.go
```

3. Move the compiled executable to any directory that is included in your PATH, for example:

```sh
mv kubectl-mergekubeconfig /usr/local/bin/
```

### Using the plugin:

To use the plugin, you can call it like other `kubectl` commands by adding the `mergekubeconfig`:

```sh
kubectl mergekubeconfig --kubeconfig1=<path_to_first_kubeconfig> --kubeconfig2=<path_to_second_kubeconfig> --output=<path_to_merged_kubeconfig>
```

