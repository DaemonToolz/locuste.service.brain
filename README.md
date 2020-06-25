# locuste.service.brain
LOCUSTE : Unité de contrôle principale

<img width="2575" alt="locuste-mcu-banner" src="https://user-images.githubusercontent.com/6602774/84285947-5a540800-ab3e-11ea-9fe2-1b9986c166b5.png">

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/4d77818f7e8b4308b2ae76b581af6c07)](https://www.codacy.com/manual/axel.maciejewski/locuste.service.brain?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=DaemonToolz/locuste.service.brain&amp;utm_campaign=Badge_Grade)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.brain&metric=alert_status)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.brain)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.brain&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.brain)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.brain&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.brain)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.brain&metric=security_rating)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.brain)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.brain&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.brain)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.brain&metric=bugs)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.brain)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=DaemonToolz_locuste.service.brain&metric=coverage)](https://sonarcloud.io/dashboard?id=DaemonToolz_locuste.service.brain)


Le project Locuste se divise en 4 grandes sections : 
* Automate (Drone Automata) PYTHON (https://github.com/DaemonToolz/locuste.drone.automata)
* Unité de contrôle (Brain) GOLANG (https://github.com/DaemonToolz/locuste.service.brain)
* Unité de planification de vol / Ordonanceur (Scheduler) GOLANG (https://github.com/DaemonToolz/locuste.service.osm)
* Interface graphique (UI) ANGULAR (https://github.com/DaemonToolz/locuste.dashboard.ui)

![Composants](https://user-images.githubusercontent.com/6602774/83644711-dcc65000-a5b1-11ea-8661-977931bb6a9c.png)

Tout le système est embarqué sur une carte Raspberry PI 4B+, Raspbian BUSTER.
* Golang 1.11.2
* Angular 9
* Python 3.7
* Dépendance forte avec la SDK OLYMPE PARROT : (https://developer.parrot.com/docs/olympe/, https://github.com/Parrot-Developers/olympe)


![Vue globale](https://user-images.githubusercontent.com/6602774/85581723-20562c00-b63d-11ea-8e0c-372c04aef6cd.png)


Détail des choix techniques pour la partie Unité de Contrôle :

* [Golang] - Rédaction rapide et simple de programmes orientés web, multithreading et multiprocessing intégré au langage
* [RPC] - Une des méthodes de communication les plus rapide
* [SocketIO] - Elément facile intégré avec Angular, Node et Python
* [ZMQ] - Système de messaging simple et rapide

Evolutions à venir : 
* Scission du serveur de socket en deux serveurs distincts afin de mieux répartir la charge (Opérateurs / Automates Python)
* Ajout de versions en GOLANG (intégration PIC)
* Intégration NGINX Reverse Proxy 
* Scission des projets en modules plus petits et partagés
