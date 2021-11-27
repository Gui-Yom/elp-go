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

L'éxecutable final se trouve dans `target/`.

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

### Fonctionnement réseau

Client :

- Connexion au serveur
- Envoi du scénario (carte + agents + tâches + paramètres)
- Récupération des résultats
- Fermeture de la connexion

Serveur (pour chaque client):

- Récupération du scénario
- Résolution du scénario
- Envoi des résultats (opérations + statistiques)
- Retour à l'étape 1

### Organisation du code

- `main.go`: Point d'entrée partagé, parsing des arguments
- `client.go`: Code relatif au client, gui
- `server.go`: Code relatif au serveur
- `net.go`: Outils réseau
- `pathfinding/`: Algorithmes relatif au pathfinding des agents
    * `astar.go`: Implémentation de A*
    * `dijkstra.go`: Implémentation de Dijkstra
    * `heuristics.go`: Implémentations de différentes fonctions heuristiques pour A*
    * `pathfinding.go`: Interfaces et utilités
    * `queue.go`: Implémentation d'une priority queue
- `scenario/`: Structures et définitions des agents, tâches, cartes, scénarios
    * `agent.go`: Agents
    * `map.go`: Cartes et tiles
    * `scenario.go`: Scenario et tâches

## Références

https://doi.org/10.1137/1.9781611973198.7
