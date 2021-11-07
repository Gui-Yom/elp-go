# elp-go

## Fonctionnement réseau

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

## Références

https://doi.org/10.1137/1.9781611973198.7
