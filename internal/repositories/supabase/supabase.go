package supabase

import (
	"github.com/go-kit/log"
	supa "github.com/nedpals/supabase-go"
)

type Supabase struct {
	logger log.Logger
	Client *supa.Client
}

func NewSupabase(logger log.Logger, client *supa.Client) *Supabase {
	return &Supabase{
		logger: logger,
		Client: client,
	}
}

func InitSupabase(supabaseUrl string, supabaseKey string) *supa.Client {
	client := supa.CreateClient(supabaseUrl, supabaseKey)
	return client
}
