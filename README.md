# elp-go

## Build

Sur Windows :

```shell
> build.bat
```

Sur Unix :

```shell
$ ./build.sh
```

Exécutable : `target/elp-go`

## Git

Pour récupérer le repo :

```shell
git clone https://github.com/Gui-Yom/elp-go
```

Pour mettre à jour le repo local :

```shell
git pull
```

Pour commit :

```shell
git add *
git commit -m "<message>"
git push origin
```

## Détails

### CLI

`$ elp-go -h`

```
Usage : elp-go [-args] [map file | rand] [width] [height] [fill] [seed]
  -addr value
        Host to connect to.
  -noconnect
        Do not connect to remote server.
  -nogui
        Disable GUI.
  -p int
        Specify the port to connect or listen to. (default 32145)
  -server
        Start a server.
Example usage :
  Start a server :
    $ elp-go -server
  Start a client with a map file :
    $ elp-go map.map
  Start a client with a randomly generated map :
    $ elp-go rand 100 100 0.1 42
```

### Organisation du code

- `main.go`: Point d'entrée partagé, parsing des arguments
- `internal/`
    - `agent.go`: Agents et tâches
    - `client.go`: Code relatif au client, gui
    - `net.go`: Outils réseau
    - `scenario.go`: Scenario et tâches
    - `server.go`: Code relatif au serveur
    - `pathfinding/`
        - `astar.go`: Implémentation de A*
        - `dijkstra.go`: Implémentation de Dijkstra
        - `heuristics.go`: Implémentations de différentes fonctions heuristiques pour A*
        - `map.go`: Cartes et tiles
        - `pathfinding.go`: Interfaces et utilités
    - `queue/`
        - `queue.go`: Interface queue
        - `linked.go`: Queue basée sur une linked list
        - `pairing.go`: Queue basée sur un pairing heap

Fichiers `*_test.go` : Code de test et benchmark

### Dépendances

- [gioui](https://gioui.org) : GUI cross-plateforme
- [giocanvas](https://github.com/ajstarks/giocanvas) : API Canvas pour gioui
- [cbor](https://github.com/fxamacker/cbor) : Serialization CBOR pour les échanges réseaux
- [testify](https://github.com/stretchr/testify) : assertions pour les tests

### Fonctionnement

Client :

- Connexion au serveur
- Envoi du scénario (carte + agents + tâches + paramètres)
- Ouverture d'une fenêtre avec la carte
- Récupération des résultats
- Fermeture de la connexion

Serveur (pour chaque client) :

- Récupération du scénario
- Résolution du scénario
- Envoi des résultats (opérations + statistiques)
- Retour à l'étape 1

## Références

https://doi.org/10.1137/1.9781611973198.7
