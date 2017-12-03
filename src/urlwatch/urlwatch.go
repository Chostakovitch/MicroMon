//urlwatch contains functions to watch URLs, i.e. check at regular intervals, using Config type as input.
//urlwatch defines the MetaResponse type, which holds websites responses' metadata.
package urlwatch

import
(
	"net/http"
)

//MetaResponse holds a website response's metadata, e.g. response code, response time, availibity, language...
type MetaResponse struct {
	Code int
	RespTime float32
	Available bool
}

func CheckUrl(url string) (MetaResponse, error) {
	//Response time = FirstByteRead - WroteRequest ! (dans httptrace) => on évite l'établissement de la connexion TLS / le DNS / ... +> c'est le vrai temps de réponse
}