Die Grundidee meiner Anwendung besteht darin, von Client zu Client in einem lokalen Netzwerk zu kommunizieren.
Die Identifizierung der Nutzer findet über IP und Name in Kombination statt, Nutzer können direkt kontaktiert werden, wenn man ihre Adresse kennt.
Sollte man die Adresse nicht kennen oder einen Nutzer erreichen wollen, der einem nicht bekannt ist, soll die Discover funktion helfen.
Geplant ist, das Clients sich gegenseitig erkennen indem sie Broadcasts in regelmäßigen Abständen im lokalen Netzwerk senden und empfangen.
Diese sollen zu einer Nutzerliste kombiniert werden.
Der Chat soll nur für lokale Netzwerke funktionieren und ist deshalb perfekt für kleiner Firmen, ohne Spuren zu hinterlassen (außer im Netzwerk)
Datenschutz und wenig bis keine Nachvollziehbarkeit (außer über die logs) sind Ziel des Konzepts. Nachrichten gehen beim neu Laden des Frontends verloren.

Im Rahmen eines Schulprojekts entstanden im Zeitraum August bis Dezember 2022


Das Backend ist in Golang verfasst und liefert ein JavaScript-Frontend aus. Diese kommunizieren über Websockets in beide Richtungen. Es wurde versucht, Multithreading in eine Art “Microservice-Architektur” zu verwandeln, zwischen den einzelnen Prozessen herrscht eine lose Kopplung über eine Thread-Sichere Datenstruktur, die den Zugriff auf die Ressource eigenständig reguliert. Die Anwendung besteht aus drei Hauptprozessen, zwei von ihnen geben ihre gesammelten Daten über zwei Websockets direkt an das Frontend weiter, um sie darzustellen. Innerhalb der empfangenden Prozesse wird je Nachricht ein Thread startet, um diese zu verarbeiten, damit der Hauptprozess nicht blockiert wird. Es laufen also mindestens zehn Threads gleichzeitig (ExplorerListener, Explorer, MessageSender, MessageListener, und innerhalb des FrontendController: je 2 Threads pro Channel.

Architekturskizze:

![Architekturskizze (Bild)](https://github.com/mortalinstrument/chatapp-client/blob/master/Diagramme%20und%20Screenshots/chatapp_architekturskizze.drawio.png?raw=true)

Klassendiagramm:

![Klassendiagramm (Bild)](https://github.com/mortalinstrument/chatapp-client/blob/master/Diagramme%20und%20Screenshots/chatapp.drawio%20(1).png?raw=true)

verwendete Libraries:
  	github.com/gorilla/websocket
  	github.com/kelseyhightower/envconfig
	  github.com/seancfoley/ipaddress-go
	  VueJS 
	  TypeScript
