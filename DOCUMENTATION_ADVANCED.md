# 📚 Documentation Technique : Project-Whis Advanced

## 1. Architecture Globale
Project-Whis Advanced est un framework de Command & Control (C2) modulaire conçu pour la furtivité et la persistance. L'architecture repose sur trois piliers :
*   **C2 Server** : Backend de gestion et interface d'administration.
*   **Agents (Linux & Windows)** : Binaires furtifs déployés sur les cibles.
*   **Post-Exploitation Tools** : Outils spécialisés pour l'évasion et le pivotement.

---

## 2. Le Serveur C2 (`/C2`)

Le serveur est le cerveau du framework. Il gère les communications, stocke les données et fournit l'interface utilisateur.

### 2.1 Backend (Go + Gorilla Mux)
Situé dans `C2/core/`, le backend assure :
*   **Routing (`Router.go`)** : Gère les endpoints de communication (API) et les pages d'administration. Il a été "Linux-ifié" pour traiter les agents Linux et Windows de manière unifiée.
*   **Database (`SQL.go`)** : Utilise MySQL pour la persistance des données (tables `windows_clients`, `linux_clients`, `commands`, `tasks`).
*   **Encryption** : Toutes les données échangées sont chiffrées en **XXTea** avec une couche **Base64** pour imiter le trafic web légitime.

### 2.2 Frontend (HTML/JS/Bootstrap)
Situé dans `C2/static/`, il comprend :
*   **Dashboards** : Visualisation en temps réel du parc d'agents.
*   **Client Management** : Pages dédiées pour Windows et Linux permettant d'envoyer des commandes spécifiques.
*   **Tasks System** : Automatisation de l'envoi de commandes basées sur des filtres (OS, Pays, Privilèges).

---

## 3. L'Agent Linux (`/Clients/HTTPS/Linux`)

L'agent Linux a été transformé pour une compatibilité universelle et une survie maximale.

### 3.1 Mécanismes de Persistance (`core/Persistence.go`)
*   **Systemd** : Création de services au niveau utilisateur (`~/.config/systemd/user/`) ou système (`/etc/systemd/system/`).
*   **SysV Init** : Support des scripts `/etc/init.d/` pour les systèmes sans Systemd (BusyBox/IoT).
*   **Crontab** : Persistance via des entrées `@reboot`.
*   **Shell Injection** : Backdoor silencieuse dans `.bashrc`, `.zshrc`, et `.profile`.
*   **SSH Key Persistence** : Injection d'une clé publique dans `authorized_keys`.
*   **Guardian Mode** : Boucle de surveillance qui redémarre l'agent s'il est tué par l'utilisateur ou un administrateur.

### 3.2 Portabilité Multi-Architecture
Grâce à une réécriture sans dépendances externes (Cgo-free), l'agent supporte :
*   **Architectures 32/64-bit** : x86, AMD64.
*   **ARM** : Idéal pour Raspberry Pi et IoT.
*   **MIPS / MIPSLE / MIPS64** : Support total pour les routeurs et équipements réseau.

### 3.3 Capacités de Communication (`core/HTTPClients.go`)
L'agent utilise un protocole de "Beaconing" : il interroge périodiquement le serveur pour récupérer des ordres, avec un délai aléatoire (**Jitter**) pour éviter la détection comportementale.

---

## 4. L'Agent Windows (`/Clients/HTTPS/Windows`)

L'agent Windows est spécialisé dans l'évasion des EDR et l'extraction de données.

### 4.1 Fonctions d'Évasion
*   **Smart UAC Bypass** : Escalade de privilèges via des techniques non-intrusives.
*   **Anti-Forensics** : Détection active de VirtualBox, VMware, Debuggers et Sandboxes d'analyse.
*   **Reflective Injection** : Capacité d'exécuter du code en mémoire sans toucher au disque.

### 4.2 Post-Exploitation
*   **Password Stealer** : Extraction des secrets de plus de 50 navigateurs.
*   **Crypto Clipper** : Détournement de transactions financières.

---

## 5. Boîte à Outils (`/Tools`)

Ces outils étendent les capacités de l'agent lors des phases de post-exploitation.

*   **Linux Fileless Injector** : Télécharge un binaire et l'exécute directement depuis la RAM via `memfd_create`. **Zéro trace sur le disque.**
*   **Socks5 Client** : Transforme une cible en tunnel réseau, permettant d'accéder aux réseaux internes de l'entreprise depuis le serveur C2.
*   **File Size Pumper** : Augmente la taille d'un binaire (ex: de 2MB à 100MB) pour saturer les scanners AV ou passer sous leur seuil d'analyse.
*   **Command Decoder** : Outil de diagnostic pour déchiffrer les flux réseaux capturés.

---

## 6. Guide de Déploiement Rapide

1.  **Serveur** : Configurez `config.toml` avec vos accès MySQL, puis lancez le binaire C2.
2.  **Compilation** :
    ```bash
    # Exemple pour un routeur MIPSLE
    cd Clients/HTTPS/Linux
    GOOS=linux GOARCH=mipsle go build -o agent_mips .
    ```
3.  **Infection** : Déployez le binaire sur la cible. Au premier lancement, il s'installera automatiquement de manière persistante.
4.  **Contrôle** : Connectez-vous à l'interface web du C2 pour piloter l'agent.

---

## 9. Intégration Havoc C2 (Architecture Recommandée)

Pour une utilisation professionnelle, il est recommandé d'utiliser **Havoc C2** comme serveur de contrôle et les agents **Project-Whis** comme implants.

### 9.1 Fonctionnement du Traducteur
Le framework inclut un script `Havoc-Integration/WhisTranslator.py`. Ce script agit comme un pont :
1.  **Agents** <--- (HTTP/XXTea) ---> **Translator** <--- (Havoc API) ---> **Havoc Teamserver**

### 9.2 Avantages
*   **Implant Custom** : Havoc ne connaît pas nativement les agents Project-Whis, ce qui les rend plus difficiles à détecter que le démon standard (Demon).
*   **Post-Exploitation Havoc** : Profitez de la puissance de Havoc (Token manipulation, Tasking, UI) tout en gardant vos méthodes de persistance Linux/MIPS.

### 9.3 Installation de l'Intégration
1.  Installez Havoc C2.
2.  Dans le dossier `Havoc-Integration`, installez la dépendance : `pip install xxtea-py`.
3.  Lancez le traducteur : `python3 WhisTranslator.py`.
4.  Configurez vos agents avec l'IP du traducteur comme URL C2.
