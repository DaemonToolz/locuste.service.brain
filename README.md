# locuste.service.brain
LOCUSTE : Unité de contrôle principale


Le project Locuste se divise en 3 grandes sections : 
* Automate (Drone Automata) PYTHON
* Unité de contrôle (Brain) GOLANG
* Unité de planification de vol / Ordonanceur (Scheduler) GOLANG
* Interface graphique (UI) ANGULAR


![Composants](https://user-images.githubusercontent.com/6602774/82243830-8960ca80-9940-11ea-917e-15585f178c6d.png)

Tout le système est embarqué sur une carte Raspberry PI 4B+, Raspbian BUSTER.
* Golang 1.11.2
* Angular 9
* Python 3.7
* Dépendance forte avec la SDK OLYMPE PARROT : (https://developer.parrot.com/docs/olympe/, https://github.com/Parrot-Developers/olympe)


![Vue globale](https://user-images.githubusercontent.com/6602774/82240232-59162d80-993a-11ea-8f8e-c7d3cfde2a7c.png)


Détail des choix techniques pour la partie Unité de Contrôle :

* [Golang] - Rédaction rapide et simple de programmes orientés web, multithreading et multiprocessing intégré au langage
* [RPC] - Une des méthodes de communication les plus rapide
* [SocketIO] - Elément facile intégré avec Angular, Node et Python

Evolutions à venir : 
* Scission du serveur de socket en deux serveurs distincts afin de mieux répartir la charge (Opérateurs / Automates Python)
* Ajout de versions en GOLANG (intégration PIC)
* Intégration NGINX Reverse Proxy 
* Scission des projets en modules plus petits et partagés
