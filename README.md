# Project-Whis Advanced C2 Framework

**Project-Whis Advanced** est un framework de Command & Control (C2) de nouvelle génération, écrit principalement en Go. Il est conçu pour la gestion à distance d'agents hautement persistants et furtifs sur des infrastructures multi-plateformes (Windows, Linux, IoT).

## 🚀 Vue d'ensemble

Ce projet offre une solution complète pour l'administration distante sécurisée, avec un serveur centralisé et des agents capables de survivre aux environnements les plus hostiles, y compris les systèmes embarqués et les distributions Linux minimalistes (BusyBox).

### Points Forts
*   **Multi-Architecture** : Support complet de x86, x64, ARM, et toutes les variantes MIPS.
*   **Persistance Hybride** : Mécanismes redondants pour garantir la survie de l'agent.
*   **Furtivité Avancée** : Anti-Forensics, Anti-VM et détection d'environnements d'analyse intégrés.
*   **Communication Chiffrée** : Flux HTTPS sécurisés avec chiffrement multiniveau (AES, XXTea).

---

## 💻 Agent Linux (Advanced)

L'agent Linux a été optimisé pour une compatibilité universelle, des serveurs Enterprise aux routeurs domestiques.

### Mécanismes de Persistance
L'agent utilise une stratégie d'installation en "Smart Copy" pour se déguiser en service système légitime.
*   **Systemd (User & System)** : Création automatique de fichiers `.service` dans `~/.config/systemd/user/` ou `/etc/systemd/system/`.
*   **SysV Init / BusyBox** : Support de `/etc/init.d/` pour les systèmes sans Systemd.
*   **Crontab Injection** : Persistance via des entrées `@reboot` pour une exécution au démarrage.
*   **Shell Profile Backdoor** : Injection silencieuse dans `.bashrc`, `.zshrc` et `.profile`.
*   **SSH Authorized Keys** : Ajout d'une clé publique pour maintenir un accès shell permanent.
*   **Guardian Mode** : Boucle de surveillance en arrière-plan qui redémarre l'agent instantanément s'il est tué.

### Architectures Supportées
L'agent est compilable pour :
*   **AMD64 / x86_64** : Serveurs et Desktops 64-bit.
*   **386 / x86** : Anciens systèmes et automates 32-bit.
*   **ARM (v5, v6, v7, ARM64)** : Raspberry Pi, IoT, Mobiles.
*   **MIPS / MIPSLE** : Routeurs, Équipements réseau.
*   **MIPS64 / MIPS64LE** : Matériel réseau haute performance.

---

## 🪟 Agent Windows (Advanced)

L'agent Windows est doté de capacités d'évasion sophistiquées.

### Fonctionnalités Clés
*   **Smart UAC Bypass** : Utilise plus de 10 méthodes différentes (FODHelper, SilentCleanup, etc.) pour obtenir les privilèges Admin.
*   **Anti-Forensics** : Détection de Debuggers, Sandbox (Any.run, VirusTotal) et environnements virtualisés.
*   **UserKit** : Dissimulation de fichiers, protection de processus critique et surveillance du registre.
*   **Stealer Intégré** : Extraction des mots de passe, cookies et cartes de crédit de plus de 50 navigateurs (Chromium & Firefox).
*   **Crypto-Clipper** : Remplacement dynamique d'adresses Bitcoin, Ethereum et Monero dans le presse-papier.

---

## 🛠️ Fonctionnalités de Commandes (C2)

Le serveur permet d'exécuter une vaste gamme de tâches sur les agents :
*   **Remote Shell** : Accès shell interactif via WebSocket.
*   **Reverse Proxy & Socks5** : Utilisation des agents comme relais réseau.
*   **DDoS Engine** : Protocoles supportés : HTTP Get, TCP Flood, UDP Flood, Slowloris, GoldenEye, SYN Flood.
*   **File Manager** : Navigation, téléchargement et exécution de fichiers à distance.
*   **Info Gathering** : Collecte détaillée de la configuration matérielle, logicielle et réseau.
*   **Torrent Seeding** : Utilisation de la bande passante pour le partage de fichiers.

## 🛠️ Outils Inclus (Tools)

Le framework est accompagné d'une suite d'outils pour faciliter le déploiement et l'évasion :
*   **Linux Fileless Injector** : Télécharge et exécute un binaire directement en RAM via `memfd_create` (furtivité totale, pas de traces sur le disque).
*   **File Size Pumper** : Augmente artificiellement la taille des binaires pour contourner les scanners antivirus.
*   **Command Decoder** : Utilitaire pour déchiffrer manuellement les communications entre les agents et le C2.
*   **Socks5 Client** : Transforme n'importe quel agent en relais Socks5 pour le pivotement réseau.

---

## 📦 Installation et Compilation

### Prérequis
*   Go 1.18+
*   Git

### Compilation de l'Agent Linux
```bash
cd Clients/HTTPS/Linux
# Pour Linux standard
go build -o agent_linux .
# Pour ARM (IoT)
GOOS=linux GOARCH=arm GOARM=7 go build -o agent_arm .
# Pour MIPS (Routeurs)
GOOS=linux GOARCH=mipsle go build -o agent_mips .
```

---

## ⚠️ Avertissement Légal (Disclaimer)

**Ce logiciel est fourni à des fins éducatives et de recherche uniquement.**
L'utilisation de cet outil pour accéder à des systèmes sans autorisation préalable est strictement interdite et illégale. Les développeurs déclinent toute responsabilité en cas de mauvaise utilisation de ce code. En téléchargeant ou en utilisant ce framework, vous acceptez l'entière responsabilité de vos actes.

---

## 🤝 Crédits et Contributions

Développé par **Bilgassim**. Basé sur les concepts originaux de Project-Whis.
Les contributions sont les bienvenues via Pull Requests pour améliorer la sécurité et la portabilité du framework.
