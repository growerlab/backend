package resolver

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Namespace() NamespaceResolver {
	return &namespaceResolver{r}
}

type mutationResolver struct{ *Resolver }
type namespaceResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
