# convertix

[![CI](https://github.com/Dhia-sel/convertix/actions/workflows/ci.yml/badge.svg)](https://github.com/Dhia-sel/convertix/actions/workflows/ci.yml)

API HTTP minimaliste de conversion d'unités, écrite en Go (bibliothèque standard uniquement), conçue comme support d'un pipeline CI/CD sécurisé : image Docker durcie, contrôles qualité automatiques, publication versionnée sur GHCR.

## Endpoints

| Méthode | Route      | Rôle                                      |
|---------|------------|-------------------------------------------|
| GET     | `/healthz` | Sonde de santé (répond `200` si le service tourne) |
| GET     | `/convert` | Convertit une valeur d'une unité vers une autre |

Paramètres de `/convert` : À COMPLÉTER (noms des paramètres de requête et exemple d'appel).

Le port d'écoute se règle avec la variable d'environnement `PORT` (défaut : `8080`).

## Lancer en local

```bash
docker build -t convertix:dev .
docker run --rm -p 8080:8080 convertix:dev
curl http://localhost:8080/healthz
```

`Ctrl+C` déclenche un arrêt gracieux : les requêtes en cours disposent de 10 secondes pour se terminer.

## Tests et qualité

```bash
bash scripts/check-coverage.sh
```

Le script lance les tests avec le détecteur de data races (`-race`), calcule la couverture totale et échoue sous le seuil (80 % par défaut, modifiable avec `COVERAGE_THRESHOLD`).

Le lint utilise golangci-lint v2 (version épinglée dans la CI) avec la configuration de `.golangci.yml`.

## CI/CD

| Workflow | Déclencheur | Rôle |
|----------|-------------|------|
| `ci.yml` | push sur `main`, pull request | lint + tests avec seuil de couverture |
| `release.yml` | tag `vX.Y.Z` | construit l'image et la publie sur GHCR (`ghcr.io/dhia-sel/convertix`) |

Les mises à jour des actions, de l'image de base et des dépendances sont proposées chaque semaine par Dependabot.

## Choix techniques

- **Image multi-stage** : compilation dans `golang:alpine`, exécution dans `distroless/static:nonroot` (aucun shell, aucun gestionnaire de paquets, utilisateur `65532`). Image d'environ 3,7 Mo compressée.
- **Binaire statique** (`CGO_ENABLED=0`, `-trimpath`, `-s -w`) : compatible distroless, plus petit, sans chemins locaux.
- **ENTRYPOINT en exec form** : le binaire est le PID 1 et reçoit directement `SIGTERM`, ce qui rend l'arrêt gracieux utile sous Kubernetes.
- **Timeouts HTTP explicites** sur le serveur, contre les connexions lentes.
- **Moindre privilège dans la CI** : `permissions: contents: read` par défaut, `packages: write` uniquement pour le job de publication.
- **Un seul fichier de version Go** : la CI lit la version dans `go.mod`, alignée avec le `Dockerfile`.

## Feuille de route

- [ ] Scan de l'image (Trivy) dans la CI
- [ ] GitOps : déploiement sur Kubernetes (kind) avec ArgoCD et Kustomize
- [ ] Déploiement AWS avec authentification OIDC (sans clé longue durée) et ECR

## Licence

Voir [LICENSE](LICENSE).