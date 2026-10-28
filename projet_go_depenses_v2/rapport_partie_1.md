# Rapport — Partie 1 : Bases de Go et serveur HTTP

Routes : `/` affiche l'accueil, `/about` décrit l'application, `/ping` répond `pong` avec GET et renvoie 405 Method Not Allowed pour les autres méthodes. `http.ServeMux` sélectionne le gestionnaire correspondant. `*http.Request` contient la méthode, l'URL, les en-têtes et le corps de la requête ; `http.ResponseWriter` permet d'écrire les en-têtes, le statut et le corps de la réponse. GET sert généralement à consulter une ressource, tandis que POST envoie des données au serveur. L'application se lance avec `go run .`.
