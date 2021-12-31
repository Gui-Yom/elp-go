# elp-go

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

## CLI

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

## Tests

CLI `go test` : https://pkg.go.dev/cmd/go#hdr-Testing_flags

```shell
$ go test -c -o target/test.<name> <package>
$ target/test.<name> -test.v
```

## Benchmarks

```shell
$ go test -c -o target/bench.pathfinding.exe elp-go/internal/pathfinding
$ target/bench.pathfinding -test.v -test.paniconexit0 -test.bench . -test.run ^$ -test.benchtime 5s -test.benchmem
```

Comparer les performances depuis un certain commit :

```shell
$ benchdiff --base-ref=d1190c07f745946dd36d8894e64562874ebf32f9 --bench=AstarLinked1000$ --benchmem --benchtime=1000x --packages=elp-go/internal/pathfinding
```

### Profilage

Il est nécessaire de limiter les benchmarks pour chaque séquence de profilage. CPU :

```shell
$ target/bench.pathfinding -test.v -test.paniconexit0 -test.bench <spec> -test.run ^$ -test.benchtime 10s -test.outputdir target -test.cpuprofile pathfinding.<spec>.cpu.prof
```

Mémoire :

```shell
$ target/bench.pathfinding -test.v -test.paniconexit0 -test.bench <spec> -test.run ^$ -test.benchtime 10s -test.benchmem -test.outputdir target -test.memprofile pathfinding.<spec>.mem.prof
```

Lecture des données de profilage avec pprof. Interface web interactive avec `pprof -http : cpu.prof`

#### Bottlenecks identifiés

Pour Dijkstra :

1. Priority queue (push O(n) pour la linked list et push/pop pour le pairing heap)
2. Map des coûts (get/put)
3. GetNeighbors (appel pour chaque tile, allocations en masse)

Pour A* (petite carte 300) :

1. Map des coûts
2. Priority queue
3. GetNeighbors

Pour A* (grande carte 10000) :

1. Map des coûts
2. Map des tiles
3. Priority queue & GetNeighbors (impact similaire)

On note aussi un coût non négligeable dû au boxing (interface{}) des ifaces.

#### Pistes d'améliorations

1. Meilleure implémentation de priority queue (linked list O(n) -_-, impl pairing heap probablement foireuse). Meilleure
   gestion de la mémoire en sachant la taille du problème afin de limiter les allocs.
2. Map optimisée Position -> Coût. Position = 2 entier 32 bits, Coût = float32. Il existe des impl de map bien mieux
   qu'une hashmap pour ce genre de données.
3. Utilisation d'un switch au lieu d'une map pour les tiles. Empeche d'avoir des tiles définies dynamiquement. Ou alors
   stocker la définition de tile directement dans la carte, multiplie par 8 la taille de la carte mais enleve toute les
   branches.
4. Spécialisation des queues pour éviter le boxing. Casse l'architecture en package du projet.
5. Précalcul des voisins pour éviter les allocs répétées au runtime. Pourrait quadrupler la taille en mémoire donc
   moyen.

## Détails

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

https://github.com/brentp/intintmap

## A faire

- faire le benchmark pour goroutine/ sans goroutine
- ajouter des tâches
- faire l'agrégation des stats dans agent.go
- récupérer les chemins pour l'interface graphique
- simulation : que les agents se déplacent
- le temps! (simulé peut-être)
- attribuer les tâches en fonction de si les petits bonhommes sont à côté ou pas
- gérer les collisions ?
