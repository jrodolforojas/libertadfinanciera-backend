package supabase

import (
	supa "github.com/nedpals/supabase-go"
)

type Supabase struct {
	Client *supa.Client
}

func NewSupabase(client *supa.Client) *Supabase {
	return &Supabase{
		Client: client,
	}
}

func InitSupabase(supabaseUrl string, supabaseKey string) *supa.Client {
	client := supa.CreateClient(supabaseUrl, supabaseKey)
	return client
}
