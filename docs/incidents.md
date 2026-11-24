# Incidents et post-mortems

Ce document consigne les pannes rencontrées pendant le projet convertix (CI/CD et GitOps). Chaque incident suit le même format : symptôme, cause racine, correction, prévention. Le but est de pouvoir expliquer chaque panne et ce qui l'empêche de revenir.

## Récapitulatif

| # | Incident | Impact | Prévention |
|---|---|---|---|
| 1 | Seuil de couverture non atteint (60 %, puis 79,5 %) | contrôle de couverture rouge | code testable, seuil de 80 % vérifié par `scripts/check-coverage.sh` |
| 2 | `gofmt` rouge sur des fichiers Go | lint rouge | `*.go text eol=lf` dans `.gitattributes` |
| 3 | Dépôt de manifests créé dans OneDrive, déplacement bloqué | travail interrompu | créer les dépôts dans `C:\dev`, hors dossier synchronisé |
| 4 | Fichiers vides, `kubectl kustomize` répond « empty » | validation impossible | contrôler les tailles avant de valider |

## 1. Seuil de couverture non atteint

**Période** : 11/11 au 20/11/2026.

**Symptôme** : le contrôle de couverture échouait. Le total était d'abord à 60 %, puis à 79,5 % le 18/11, sous le seuil de 80 %.

**Cause racine** : du code n'était pas testé. Le package `main` n'avait aucun test, et `main()` mélangeait la lecture de la configuration et le démarrage du serveur, ce qui le rendait impossible à tester. Le converter, lui, n'avait pas de défaut de code : il manquait seulement des tests.

**Correction** : en trois commits.
- `b07dc46` (11/11) : tests du converter (conversions et erreurs).
- `8f6b9eb` (18/11) : ajout de `main_test.go`.
- `df181f7` (20/11) : `listenAddr()` et `run()` extraits de `main()` avec un cycle de vie testable, et `TestListenAddr` avec `t.Setenv`.

La couverture totale est passée à 83,8 %, sans modifier le code du converter.

**Prévention** : le seuil de 80 % est vérifié par `scripts/check-coverage.sh`. La vraie leçon est de rendre le code testable en séparant la logique du démarrage, plutôt que d'écrire des tests artificiels pour monter un pourcentage. Les messages de commit doivent aussi décrire le changement réel, pour retrouver l'historique sans se fier à la mémoire.

## 2. `gofmt` rouge sur des fichiers Go

**Période** : 19 au 20/11/2026.

**Symptôme** : `gofmt` signalait `main.go`, `main_test.go` et `handler_test.go` comme mal formatés, alors que le code était correct.

**Cause racine** : les fins de ligne. Les fichiers créés sous Windows utilisaient CRLF, alors que `gofmt` attend du LF.

**Correction** : règle `*.go text eol=lf` ajoutée dans `.gitattributes` (commit `df181f7`), puis `golangci-lint fmt ./...`.

**Prévention** : le `.gitattributes` est versionné, donc la règle s'applique à tout nouveau clone. Les scripts `*.sh` étaient déjà forcés en LF.

## 3. Dépôt de manifests créé dans OneDrive, déplacement bloqué

**Date** : 23/11/2026.

**Symptôme** : `Move-Item` du dépôt `convertix-manifests` vers `C:\dev` échouait avec l'erreur « Le processus ne peut pas accéder au fichier, car il est utilisé par un autre processus ».

**Cause racine** : deux problèmes distincts.
- Le dépôt avait été créé sur le Bureau, un dossier synchronisé par OneDrive, alors que la règle du projet est de travailler dans `C:\dev`.
- Des processus VS Code tenaient des fichiers du dossier : `Get-Process code` en montrait 15, puis 10 après une première fermeture, avec une seule fenêtre « Welcome » encore ouverte.

**Correction** :
1. Sortir du dossier (`cd C:\`).
2. Arrêter VS Code (`Stop-Process -Name code -Force`) après avoir enregistré.
3. Relancer `Move-Item`.
4. Vérifier : ancien chemin absent (`Test-Path` à `False`), `git status` et `git remote -v` inchangés, tailles des fichiers identiques.

**Prévention** : créer chaque dépôt directement dans `C:\dev`. OneDrive synchronise aussi le dossier `.git`, ce qui peut verrouiller des fichiers pendant une opération Git ou créer des copies en conflit.

## 4. Fichiers vides et validation Kustomize

**Date** : 23/11/2026.

**Symptôme** : `kubectl kustomize overlays\dev` répondait « empty ».

**Cause racine** : le fichier `overlays/dev/kustomization.yaml` n'était pas encore écrit (0 octet). C'était le cas de 5 fichiers du dépôt à ce moment-là : `README.md`, `argocd/application-dev.yaml`, `base/kustomization.yaml`, `overlays/dev/kustomization.yaml` et `overlays/prod/kustomization.yaml`.

**Correction** : écrire le fichier, l'enregistrer (`Ctrl+S`), puis relancer `kubectl kustomize overlays\dev`.

**Prévention** : avant toute validation, lister les tailles de tous les fichiers :

```powershell
Get-ChildItem -Recurse -File | Where-Object { $_.FullName -notmatch '\\\.git\\' } | Select-Object FullName, Length
```

Un fichier à `0` signifie qu'il n'est pas écrit ou pas enregistré. Le contrôle prend une seconde et évite de chercher une erreur de YAML qui n'existe pas.