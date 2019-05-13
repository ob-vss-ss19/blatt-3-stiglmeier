## Abgabe

-   In der IDE starten 

oder

-   Mit Docker
    1. Docker images herunterladen: 
        ``` 
        docker pull terraform.cs.hm.edu:5043/ob-vss-ss19-blatt-3-stiglmeier:develop-treeservice
        docker pull terraform.cs.hm.edu:5043/ob-vss-ss19-blatt-3-stiglmeier:develop-treecli
        ```    
    2. Docker netwerk erstellen
        ```
        docker network create actors
        ```
    3. Docker images starten
        ```
        docker run --rm --net actors --name treeservice terraform.cs.hm.edu:5043/ob-vss-ss19-blatt-3-stiglmeier:develop-treeservice --bind="treeservice.actors:8091"
        docker run --rm --net actors --name treecli terraform.cs.hm.edu:5043/ob-vss-ss19-blatt-3-stiglmeier:develop-treecli --bind="treecli.actors:8090 --remote="treeservice.actors:8091"
        ```




-----------------------
## Ausführen mit Docker

-   Images bauen

    ```
    make docker
    ```

-   ein (Docker)-Netzwerk `actors` erzeugen

    ```
    docker network create actors
    ```

-   Starten des Tree-Services und binden an den Port 8090 des Containers mit dem DNS-Namen
    `treeservice` (entspricht dem Argument von `--name`) im Netzwerk `actors`:

    ```
    docker run --rm --net actors --name treeservice treeservice \
      --bind="treeservice.actors:8090"
    ```

    Damit das funktioniert, müssen Sie folgendes erst im Tree-Service implementieren:

    -   die `main` verarbeitet Kommandozeilenflags und
    -   der Remote-Actor nutzt den Wert des Flags
    -   wenn Sie einen anderen Port als `8090` benutzen wollen,
        müssen Sie das auch im Dockerfile ändern (`EXPOSE...`)

-   Starten des Tree-CLI, Binden an `treecli.actors:8091` und nutzen des Services unter
    dem Namen und Port `treeservice.actors:8090`:

    ```
    docker run --rm --net actors --name treecli treecli --bind="treecli.actors:8091" \
      --remote="treeservice.actors:8090" trees
    ```

    Hier sind wieder die beiden Flags `--bind` und `--remote` beliebig gewählt und
    in der Datei `treeservice/main.go` implementiert. `trees` ist ein weiteres
    Kommandozeilenargument, dass z.B. eine Liste aller Tree-Ids anzeigen soll.

    Zum Ausprobieren können Sie den Service dann laufen lassen. Das CLI soll ja jedes
    Mal nur einen Befehl abarbeiten und wird dann neu gestartet.

-   Zum Beenden, killen Sie einfach den Tree-Service-Container mit `Ctrl-C` und löschen
    Sie das Netzwerk mit

    ```
    docker network rm actors
    ```

## Ausführen mit Docker ohne vorher die Docker-Images zu bauen

Nach einem Commit baut der Jenkins, wenn alles durch gelaufen ist, die beiden
Docker-Images. Sie können diese dann mit `docker pull` herunter laden. Schauen Sie für die
genaue Bezeichnung in die Consolenausgabe des Jenkins-Jobs.

Wenn Sie die Imagenamen oben (`treeservice` und `treecli`) durch die Namen aus der
Registry ersetzen, können Sie Ihre Lösung mit den selben Kommandos wie oben beschrieben,
ausprobieren.
