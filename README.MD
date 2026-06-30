![image](./assets/wacdo-banner-011f2cee-9131-498e-8009-a01014ce4d34.png)

# Wacdo — Back-office de borne de commande

## Contexte

Développement du back-office d'une application de borne de commande **Wacdo**.  
L'objectif est de gérer de manière sécurisée les données de l'application, les commandes clients et les utilisateurs selon leurs rôles.

---

## Stack technique

- **Langage** : Go (Golang)
- **Base de données** : Maria DB

---

## Fonctionnalités

### Gestion des utilisateurs

Trois rôles distincts sont à implémenter :

| Rôle | Permissions |
|---|---|
| **Administration** | Gestion complète des données et des utilisateurs |
| **Préparation** | Consultation et validation des commandes |
| **Accueil** | Saisie de commandes (comptoir/téléphone) et remise au client |

### Gestion des produits et menus

> Accessible uniquement aux utilisateurs avec le rôle **Administration**.

- CRUD sur les produits : nom, description, prix, image, disponibilité
- CRUD sur les menus : composition, options disponibles

### Saisie de commandes

> Accessible aux rôles **Accueil** et **Administration**.

- Saisie d'une commande (comptoir ou centre d'appel)
- Remise d'une commande à un client → statut `livrée`
- Pas de gestion de paiement : les commandes sont identifiées par un numéro d'identification

### Préparation de commandes

> Accessible au rôle **Préparation**.

- Liste des commandes à préparer, triées par heure de livraison croissante
- Marquage d'une commande comme `préparée`

---

## API publique (communication avec le front-end)

| Méthode | Description |
|---|---|
| `GET /menus` | Liste détaillée des menus |
| `GET /products` | Liste des produits (filtrable par catégorie) |
| `POST /orders` | Réception du détail d'une commande |

---

## Sécurité

- Authentification par identifiant unique + mot de passe
- Gestion de sessions / tokens (JWT ou équivalent)
- Contrôle d'accès basé sur les rôles (RBAC)
- Protection contre les injections SQL
- Protection des données sensibles

---

## Livrables

- [ ] Schéma conceptuel (MCD) et physique (MPD) de la base de données
- [ ] Schémas fonctionnels de l'application
- [ ] Base de données opérationnelle
- [ ] Application déployée sur un serveur
- [ ] Sources et documentation sur ce dépôt GitHub

---

## Critères d'évaluation

### Base de données
- [ ] Les données nécessaires à l'application sont correctement identifiées
- [ ] Les données sont retranscrites sur un schéma décrivant les différentes tables et leurs relations
- [ ] Le nommage des tables et des champs est cohérent avec la typologie des données
- [ ] Le type des champs est choisi en adéquation avec la nature des données
- [ ] La mise en relation des tables est correctement effectuée
- [ ] Les requêtes sont affinées avec des systèmes de tri et de filtres
- [ ] Les requêtes sont optimisées par l'utilisation de clés étrangères et de jointures

### Protection des données
- [ ] Les données sensibles et réglementées bénéficient d'un traitement spécifique et sont protégées
- [ ] L'application informe l'utilisateur du stockage, de l'utilisation et du cadre de partage de ses données personnelles
- [ ] L'utilisateur dispose d'un droit de consultation, modification et de suppression de ses données personnelles
<!-- https://docs.github.com/en/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax#alerts -->
> [!NOTE]
> Cette application n'utilise pas de données personnelles . Les mots de passe sont stockées chiffrés en base de données.
### Conception fonctionnelle
- [ ] Toutes les fonctionnalités nécessaires au bon fonctionnement de l'application sont correctement listées et détaillées
- [ ] Le schéma fonctionnel décrit en détail l'enchaînement des vues en fonction des différentes actions et interactions

### Qualité du code
- [ ] Le code est indenté et les commentaires aident à la compréhension
- [ ] Les dossiers et fichiers du projet sont organisés
- [ ] Les conventions de nommage sont respectées pour l'ensemble du code
- [ ] Les limites du code sont connues et documentées
- [ ] Les erreurs de codage sont traitées

### Architecture
- [ ] Le modèle gère les interactions avec la base de données
- [ ] Les contrôleurs implémentent la logique et préparent les données nécessaires
- [ ] La vue permet l'affichage des données transmises par le contrôleur
- [ ] Le programme protège l'intégrité des données en empêchant les injections

### Sécurité & authentification
- [ ] Un utilisateur s'authentifie via un identifiant unique et un mot de passe
- [ ] Un système de session ou de token permet d'identifier l'utilisateur connecté
- [ ] Différents rôles sont implémentés et délimitent les actions possibles pour chaque type d'utilisateur

### Collaboration & déploiement
- [ ] L'utilisation de Git / GitHub est maîtrisée (commits clairs, branches, etc.)
- [ ] Les fonctionnalités déployées sont conformes au cahier des charges
- [ ] Des tests unitaires sont réalisés et validés
- [ ] L'application mise en ligne est exempte de bugs et fonctionnelle
- [ ] L'application est testée en production sans erreurs ni effets de bord

# Limite / Pour aller plus loin
* Les sessions invités / customer ne sont pas implémentées : un client ne peut pas s'enregistrer, et voir ses commandes passées, il n'y a pas de Guest Session (qui pourraient par exemple passer une commande aux bornes)
* Les modèles de ressources pourraient faire de l'embedding pour ne pas répéter les mêmes déclarations
* La gestion d'utilisateurs est "sécurisée" en autorisant l'accès uniquement depuis localhost à défaut de pouvoir implémenter un système d'IAM complet. 
* Gorm ne supporte pas la création de Menus ou Products qui sont Available (bool)=false,car les modèles en structs ne prennent pas les champs nuls, 0 ou false par défaut. On considère pour l'instant que la création ne peut se faire que en true, et qu'une update est possible (par le client) en false ensuite, car update utilise une ```map[string]any```. Une implémentation front-end pourrait le gérer (pas idéal) ou le handler http lui-même avec un pointer, qui complexifie le code. 
* Images non implémentées.

# Schéma fonctionnels

## Cas d'usage
```mermaid
flowchart LR
    Borne(["Borne / Front public<br/>non authentifie"])
    Accueil(["Accueil"])
    Prepa(["Preparation"])
    Admin(["Administration"])

    subgraph Consultation["Consultation publique"]
        P1["Consulter les produits<br/>GET /products"]
        P2["Consulter les menus<br/>GET /menus"]
    end
    subgraph Commandes["Gestion des commandes"]
        C1["Creer une commande<br/>POST /orders"]
        C2["Lister / voir les commandes<br/>GET /orders"]
        C3["Faire evoluer l'etat<br/>PUT /orders/:id"]
        C4["Supprimer une commande<br/>DELETE /orders/:id"]
    end
    subgraph Catalogue["Catalogue - Admin"]
        K1["CRUD produits"]
        K2["CRUD menus"]
    end
    subgraph Users["Comptes"]
        U1["S'inscrire / se connecter"]
        U2["Assigner un role<br/>PUT /users/:username/role"]
    end

    Borne --> P1
    Borne --> P2

    Accueil --> C1
    Accueil --> C2
    Accueil --> C3

    Prepa --> C2
    Prepa --> C3

    Admin --> C4
    Admin --> K1
    Admin --> K2
    Admin --> U2

    U1 -.-> Accueil
    U1 -.-> Prepa
    U1 -.-> Admin
```

## Cycle de vie des commandes
```mermaid
stateDiagram-v2
    [*] --> Created: Accueil cree la commande (POST /orders)
    Created --> Validated: Accueil confirme le paiement
    Validated --> Ready: Preparation marque preparee
    Ready --> Delivered: Accueil remet au client
    Delivered --> [*]
    Created --> [*]: Admin supprime (etat Created uniquement)
```

## Parcours de commande
```mermaid
flowchart TD
    A["Client choisit produits / menus<br/>sur la borne"] --> B["Accueil saisit la commande<br/>POST /orders"]
    B --> C{"Produits / menus<br/>disponibles ?"}
    C -- Non --> C1["Rejet : article indisponible"]
    C1 --> A
    C -- Oui --> D["Commande creee - etat = Created<br/>prix calcule cote serveur"]
    D --> E{"Paiement<br/>confirme ?"}
    E -- Non --> E1["Reste Created<br/>supprimable par Admin"]
    E -- Oui --> F["Accueil : etat -> Validated"]
    F --> G["Preparation prepare<br/>etat -> Ready"]
    G --> H["Accueil remet au client<br/>etat -> Delivered"]
    H --> I(["Commande terminee"])

```

## Authentification
```mermaid
sequenceDiagram
    participant U as Utilisateur (staff)
    participant API as API Wacdo
    participant DB as MariaDB

    U->>API: POST /users/login (username, password)
    API->>DB: SELECT user WHERE username
    DB-->>API: user + hash bcrypt
    API->>API: bcrypt.Compare + verifie role != nil
    API-->>U: 200 + JWT (HS256, expire 8h)

    Note over U,API: Requete sur une route protegee
    U->>API: Requete + Authorization: Bearer JWT<br/>(GET / POST / PUT / DELETE)
    API->>API: Authenticate (signature + expiration)
    API->>API: Authorize (role autorise ?)
    alt Token valide ET role autorise
        Note over API,DB: Pas de cache : la base est toujours interrogee
        API->>DB: GET -> SELECT<br/>POST -> INSERT<br/>PUT -> UPDATE<br/>DELETE -> DELETE (soft delete)
        DB-->>API: Resultat / lignes affectees
        API-->>U: 200 / 201 / 202 / 204 + donnees
    else Token absent / invalide
        API-->>U: 401 Unauthorized
    else Role insuffisant
        API-->>U: 403 / 401 Forbidden
    end
    
```
