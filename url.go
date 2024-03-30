package url

import (
	"context"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/api/option"
)

type Config struct {
	ProjectID         string `envconfig:"project_ID"`
	PrivateKeyId      string `envconfig:"private_key_id"`
	PrivateKey        string `envconfig:"private_key"`
	ClientEmail       string `envconfig:"client_email"`
	ClientId          string `envconfig:"client_id"`
	ClientX509CertUrl string `envconfig:"client_x509_cert_url"`
}

func createClient(ctx context.Context, cfg Config) *firestore.Client {
	jsonToken := []byte(`{
		"type": "service_account",
		"project_id": "` + cfg.ProjectID + `",
		"private_key_id": "` + cfg.PrivateKeyId + `",
		"private_key": "` + cfg.PrivateKey + `",
		"client_email": "` + cfg.ClientEmail + `",
		"client_id": "` + cfg.ClientId + `",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "` + cfg.ClientX509CertUrl + `",
		"universe_domain": "googleapis.com"
	  }
	  `)

	client, err := firestore.NewClientWithDatabase(ctx, cfg.ProjectID, "tiny-url-db", option.WithCredentialsJSON(jsonToken))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// Close client when done with
	// defer client.Close()
	return client
}

type URLShortener struct {
	urls map[string]string
}

func (us *URLShortener) SearchForMatch(url string) string {
	if val, ok := us.urls[url]; ok {
		return val
	}
	return "NotFound"
}

func (us *URLShortener) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		log.Printf("took %v\n", time.Since(start))
	}()
	if r.URL.Path == "/NotFound" {
		return
	}

	http.Redirect(w, r, us.SearchForMatch(r.URL.Path), http.StatusSeeOther)
	log.Println(r.URL.Path)

}

func init() {
	var cfg = Config{}
	err := envconfig.Process("App", &cfg)
	log.Println(cfg)
	if err != nil {
		log.Fatalf("Failed to process env var: %v", err)
	}

	mux := http.NewServeMux()
	client := createClient(context.Background(), cfg)

	defer client.Close()
	us := &URLShortener{}

	urls, err := client.Collection("urls").Doc("kbj2w5cHLZzSiLJFwvqj").Get(context.Background())
	if err != nil {
		log.Fatalf("Failed to get urls: %v", err)
	}
	us.urls = make(map[string]string)
	temp := make(map[string]map[string]string)
	urls.DataTo(&temp)
	us.urls = temp["urls"]
	log.Println(us.urls)

	log.Println(len(us.urls))

	mux.HandleFunc("/", us.RedirectHandler)
	functions.HTTP("url", mux.ServeHTTP)
}
