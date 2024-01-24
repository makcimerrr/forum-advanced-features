<article>
      <h1>Organisation d'un projet Go standard</h1>
<blockquote>
<p>Cet article est une traduction du dépôt Github <a href="https://github.com/golang-standards/project-layout">Standard Go Project Layout</a></p>
</blockquote>
<h2>Introduction</h2>
<p>Ce dépôt est une architecture basique pour des projets d'applications en Go. Il ne représente pas un standard officiel défini par l'équipe de développement principale de Go. C'est néanmoins un ensemble de modèles d'architecture que l'on peut retrouver autant sur des projets historiques que plus récents dans l'écosystème Go. Certains patterns sont plus populaires que d'autres. Il comporte également nombre d'améliorations mineures ainsi que plusieurs répertoires communs à beaucoup d'applications existantes de taille importante.</p>
<p>Si vous commencez à apprendre Go, ou si vous souhaitez développer un petit side-project pour vous-même, cette architecture n'est pas du tout adaptée. Commencez par quelque chose de très simple (un unique fichier <code>main.go</code> est largement suffisant). Au fur et à mesure que votre projet évolue, il est important de garder à l'esprit que votre code doit être bien structuré sinon vous finirez rapidement avec un code difficile à maintenir, comprenant beaucoup de dépendances cachées et un state global. Plus il y aura de gens qui travailleront sur le projet, plus il sera important d'avoir une structure solide. C'est pourquoi il est important d'introduire une façon identique pour tout le monde de gérer les bibliothèques et les packages. Lorsque vous maintenez un projet open source ou que vous savez que d'autres projets importent votre code depuis votre dépôt, il est important d'avoir des packages et du code privé (aka <code>internal</code>). Clonez le dépôt, gardez ce dont vous avez besoin et supprimez tout le reste ! Ce n'est pas parce que des dossiers existent que vous devez impérativement tous les utiliser. Tous ces patterns ne sont pas tout le temps utilisés dans tous les projets. Même le pattern <code>vendor</code> n'est pas universel.</p>
<p>Depuis la sortie de Go 1.14 les <a href="https://github.com/golang/go/wiki/Modules"><code>Go Modules</code></a> sont enfin prêts à être utilisés en production. Utilisez les <a href="https://blog.golang.org/using-go-modules"><code>Go Modules</code></a> par défaut sauf si vous avez une raison bien spécifique de ne pas les utiliser. Lorsque vous les utilisez, vous n'avez pas besoin de vous embêter avec le $GOPATH ou de définir le dossier dans lequel vous allez mettre votre projet. Le fichier <code>go.mod</code> part du principe que votre dépôt est hébergé sur Github, mais ce n'est pas une obligation. Le chemin du module peut être n'importe quoi, mais il faut savoir que le premier composant du chemin devrait toujours avoir un point dans son nom (la version actuelle de Go ne l'impose plus, mais si vous utilisez des versions un peu plus anciennes ne soyez pas surpris que votre build échoue s'il n'y a pas de point). Allez voir les tickets <a href="https://github.com/golang/go/issues/37554"><code>37554</code></a> et <a href="https://github.com/golang/go/issues/32819"><code>32819</code></a> si vous souhaitez en savoir plus.</p>
<p>L'architecture de ce projet est générique de manière intentionelle et elle n'essaie pas d'imposer une structure de paquet Go spécifique.</p>
<p>Ce projet est un effort communautaire. Ouvrez un ticket si vous découvrez un nouveau pattern ou si vous pensez qu'un des patterns existants devrait être mis à jour.</p>
<p>Si vous avez besoin d'aide pour le nommage, le formattage ou le style, commencez par lancer <a href="https://golang.org/cmd/gofmt/"><code>gofmt</code></a> et <a href="https://github.com/golang/lint"><code>golint</code></a>. Prenez également le temps de parcourir ces lignes directrices et recommandations :</p>
<ul>
<li><a href="https://talks.golang.org/2014/names.slide">https://talks.golang.org/2014/names.slide</a></li>
<li><a href="https://golang.org/doc/effective_go.html#names">https://golang.org/doc/effective_go.html#names</a></li>
<li><a href="https://blog.golang.org/package-names">https://blog.golang.org/package-names</a></li>
<li><a href="https://github.com/golang/go/wiki/CodeReviewComments">https://github.com/golang/go/wiki/CodeReviewComments</a></li>
<li><a href="https://rakyll.org/style-packages">Style guideline for Go packages</a> (rakyll/JBD)</li>
</ul>
<p>Lisez l'article <a href="https://medium.com/golang-learn/go-project-layout-e5213cdcfaa2"><code>Go Project Layout</code></a> pour avoir des informations additionnelles.</p>
<p>Plus d'infos sur le nommage et l'organisation des packages, ainsi que quelques recommandations sur la structuration du code :</p>
<ul>
<li><a href="https://www.youtube.com/watch?v=PTE4VJIdHPg">GopherCon EU 2018: Peter Bourgon - Best Practices for Industrial Programming</a></li>
<li><a href="https://www.youtube.com/watch?v=MzTcsI6tn-0">GopherCon Russia 2018: Ashley McNamara + Brian Ketelsen - Go best practices.</a></li>
<li><a href="https://www.youtube.com/watch?v=ltqV6pDKZD8">GopherCon 2017: Edward Muller - Go Anti-Patterns</a></li>
<li><a href="https://www.youtube.com/watch?v=oL6JBUk6tj0">GopherCon 2018: Kat Zien - How Do You Structure Your Go Apps</a></li>
</ul>
<p>Un article en Chinois sur les guidelines du design orienté package et la couche Architecture :</p>
<ul>
<li><a href="https://github.com/danceyoung/paper-code/blob/master/package-oriented-design/packageorienteddesign.md">面向包的设计和架构分层</a></li>
</ul>
<h2>Les répertoires Go</h2>
<h3><code>/cmd</code></h3>
<p>Les applications principales de ce projet.</p>
<p>Le nom de répertoire de chaque application doit correspondre au nom de l'exécutable que vous souhaitez avoir (p. ex., <code>/cmd/myapp</code>).</p>
<p>Ne mettez pas trop de code dans le répertoire de votre application. Si vous pensez que le code peut être importé et réutilisé dans d'autres projets, déplacez le dans le dossier <code>/pkg</code>. Si le code n'est pas réutilisable, ou si vous ne souhaitez pas que d'autres personnes l'utilisent, placez le dans le dossier <code>/internal</code>. Soyez explicite quant à vos intentions, vous seriez surpris de l'utilisation que d'autres développeurs pourraient faire de votre code !</p>
<p>Il est habituel d'avoir une petite fonction <code>main</code> qui importe et appelle du code contenu dans les dossiers <code>/internal</code> et <code>/pkg</code>, et rien de plus.</p>
<p>Voir le dossier <a href="https://github.com/golang-standards/project-layout/tree/master/cmd/README.md"><code>/cmd</code></a> pour des exemples.</p>
<h3><code>/internal</code></h3>
<p>Applications privées et bibliothèques de code. C'est le code que vous ne souhaitez pas voir importé dans d'autres applications ou bibliothèques. Notez que ce pattern est imposé par le compilateur Go lui-même (voir les <a href="https://golang.org/doc/go1.4#internalpackages"><code>release notes</code></a> de Go 1.4 pour plus de détails). Vous n'êtes pas limité à un seul dossier <code>internal</code> de haut niveau, mais vous pouvez en avoir plusieurs à n'importe quel niveau de l'arborescence de votre projet.</p>
<p>Vous pouvez également ajouter un peu de structure dans vos packages internes pour séparer le code partagé et non partagé. Ce n'est pas du tout obligatoire (surtout pour les petits projets), mais il est intéressant d'avoir des indices visuels indiquant l'utilisation prévue d'un package. Le code de votre application peut aller dans un dossier <code>/internal/app</code> (p. ex., <code>/internal/app/myapp</code>) tandis que le code partagé par les applications peut se retrouver dans un dossier <code>/internal/pkg</code> (p. ex., <code>/internal/pkg/myprivlib</code>).</p>
<h3><code>/pkg</code></h3>
<p>Placez-y le code qui peut être réutilisé par les applications externes (p. ex., <code>/pkg/mypubliclib</code>). D'autres projets peuvent importer ces bibliothèques et s'attendent donc à ce qu'elles soient fonctionnelles, pensez y donc à deux fois avant de mettre du code dans ce dossier :-) Utiliser le dossier <code>internal</code> est une manière plus adéquate de garder vos packages privés et non importables car c'est intégré au compilateur Go. Le dossier <code>/pkg</code> est nénanmoins une bonne manière d'indiquer que le code contenu dans ce dossier peut être utilisé par les autres utilisateurs sans problème. L'article de blog de Travis Jeffery <a href="https://travisjeffery.com/b/2019/11/i-ll-take-pkg-over-internal/"><code>I'll take pkg over internal</code></a> explique plus en détail les différences entre les dossier <code>pkg</code> et <code>internal</code> et quand il fait sens de les utiliser.</p>
<p>C'est également une manière de regrouper tout votre code Go au même endroit lorsque votre dossier racine comporte de nombreux composants et dossiers non-Go, permettant plus facilement de lancer les différents outils Go, tel que mentionné dans les conférences suivantes : <a href="https://www.youtube.com/watch?v=PTE4VJIdHPg"><code>Best Practices for Industrial Programming</code></a> lors de GopherCon EU 2018, <a href="https://www.youtube.com/watch?v=oL6JBUk6tj0">GopherCon 2018: Kat Zien - How Do You Structure Your Go Apps</a> et <a href="https://www.youtube.com/watch?v=3gQa1LWwuzk">GoLab 2018 - Massimiliano Pippi - Project layout patterns in Go</a>).</p>
<p>Voir le dossier <a href="https://github.com/golang-standards/project-layout/tree/master/pkg/README.md"><code>/pkg</code></a> pour découvrir quels projets Go populaires utilisent cette architecture de projet. C'est un pattern plutôt commun, mais qui n'est pas accepté de manière universelle, et certaines personnes de la communauté Go ne le recommandent pas.</p>
<p>Vous n'êtes pas obligés de l'utiliser si votre projet est petit et si l'ajout d'un niveau de plus n'ajoute pas vraiment de valeur (sauf si vous y tenez vraiment :-)). Il est temps d'y penser lorsque votre projet commence à prendre de l'ampleur et que votre dossier racine est encombré (surtout si vous avez beaucoup de composants non-Go)</p>
<h3><code>/vendor</code></h3>
<p>Les dépendances de votre application (gérées manuellement ou via votre gestionnaire de dépendances favori tel que la fonctionnalité incluse dans les <a href="https://github.com/golang/go/wiki/Modules"><code>Go Modules</code></a>). La commande <code>go mod vendor</code> créera un dossier <code>/vendor</code> pour vous. Notez que vous devrez peut-être utiliser le flag <code>-mod=vendor</code> avec votre commande <code>go build</code> si vous n'utilisez pas Go 1.14 qui le définit par défaut.</p>
<p>Ne commitez pas vos dépendances si vous développez une bibliothèque.</p>
<p>Depuis sa version <a href="https://golang.org/doc/go1.13#modules"><code>1.13</code></a>, Go active la fonctionnalité de proxy de module (en utilisant <a href="https://proxy.golang.org"><code>https://proxy.golang.org</code></a> comme serveur de proxy par défaut). Plus d'infos <a href="https://blog.golang.org/module-mirror-launch"><code>ici</code></a> afin de définir si cela correspond à votre obligations et contraintes. Si c'est le cas, vous n'aurez pas besoin du dossier <code>vendor</code>.</p>
<h2>Les répertoires d'application de services</h2>
<h3><code>/api</code></h3>
<p>Spécifications OpenAPI/Swagger, fichiers de schémas JSON, fichiers de définitions de protocoles.</p>
<p>Voir le dossier <a href="https://github.com/golang-standards/project-layout/tree/master/api/README.md"><code>/api</code></a> pour des examples.</p>
<h2>Les répertoires d'application web</h2>
<h3><code>/web</code></h3>
<p>Les composants spécifiques aux applications web : assets statiques, templates serveurs et SPAs.</p>
<h2>Les répertoire communs aux applications</h2>
<h3><code>/configs</code></h3>
<p>Templates de fichiers de configuration ou configurations par défaut.</p>
<p>Ajoutez vos templates <code>confd</code> ou <code>consul-template</code> dans ce répertoire.</p>
<h3><code>/init</code></h3>
<p>Initialisation du système (systemd, upstart, sysv) et configurations des administrateurs/superviseurs de process (runit, supervisord).</p>
<h3><code>/scripts</code></h3>
<p>Scripts permettant différentes opérations telles que le build, l'installation, des analyses, ...</p>
<p>Ces scripts permettent de garder le Makefile du dossier racine réduit et simple (p. ex., <a href="https://github.com/hashicorp/terraform/blob/master/Makefile"><code>https://github.com/hashicorp/terraform/blob/master/Makefile</code></a>).</p>
<p>Voir le dossier <a href="https://github.com/golang-standards/project-layout/tree/master/scripts/README.md"><code>/scripts</code></a> pour des exemples.</p>
<h3><code>/build</code></h3>
<p>Packaging et Intégration Continue.</p>
<p>Ajoutez vos scripts et configurations de cloud (AMI), conteneur (Docker), OS (deb, rpm, pkg) et package dans le dossier <code>/build/package</code>.</p>
<p>Placez vos scripts et configurations de CI (travis, circle, drone) dans le dossier <code>/build/ci</code>. Prenez garde au fait que certains outils de CI (p. ex., Travis CI) sont très contraignants vis à vis de l'emplacement de leurs fichiers de configuration. Essayez donc, lorsque c'est possible, de créer des liens entre le dossier <code>/build/ci</code> et les endroits où les outils de CI s'attendent à trouver ces fichiers.</p>
<h3><code>/deployments</code></h3>
<p>Templates et configurations pour les IaaS, PaaS, système et l'orchestration de conteneurs (docker-compose, kubernetes/helm, mesos, terraform, bosh). Sur certains projets (principalement les applications déployées via Kubernetes) ce dossier s'appelle <code>/deploy</code>.</p>
<h3><code>/test</code></h3>
<p>Applications et données de tests externes additionnels. Vous pouvez structurer le dossier <code>/test</code> de la façon qui vous convient le mieux. Pour des projets plus importants, il fait sens d'utiliser un sous-dossier <code>data</code>. Vous pouvez par exemple utiliser <code>/test/data</code> ou <code>/test/testdata</code> si vous souhaitez que Go ignore ce dossier. Go ignore également les dossiers ou fichiers commençant par "." ou "_", ce qui vous donne plus de flexibilité en terme de nommage de votre dossier de données de test.</p>
<p>Voir le dossier <a href="https://github.com/golang-standards/project-layout/tree/master/test/README.md"><code>/test</code></a> pour des exemples</p>
<h2>Autres répertoires</h2>
<h3><code>/docs</code></h3>
<p>Documents utilisateurs et design (en plus de votre documentation générée GoDoc)</p>
<p>Voir le dossier <a href="https://github.com/golang-standards/project-layout/tree/master/docs/README.md"><code>/docs</code></a> pour des exemples</p>
<h3><code>/tools</code></h3>
<p>Outils de support du projet. Ces scripts peuvent importer du code des dossier <code>/pkg</code> et <code>/internal</code></p>
<p>Voir le dossier <a href="https://github.com/golang-standards/project-layout/tree/master/tools/README.md"><code>/tools</code></a> pour des exemples</p>
<h3><code>/examples</code></h3>
<p>Exemples de vos applications et/ou de vos bibliothèques publiques</p>
<p>Voir le dossier <a href="https://github.com/golang-standards/project-layout/tree/master/examples/README.md"><code>/examples</code></a> pour des exemples</p>
<h3><code>/third_party</code></h3>
<p>Outils d'aide externe, code forké et autres utilitaires tierces (p. ex., Swagger UI).</p>
<h3><code>/githooks</code></h3>
<p>Hooks Git.</p>
<h3><code>/assets</code></h3>
<p>D'autres assets qui sont utilisés dans votre dépôt (images, logos, etc).</p>
<h3><code>/website</code></h3>
<p>C'est là que vous placez les données du site de votre projet si vous n'utilisez pas GitHub pages.</p>
<p>Voir le dossier <a href="https://github.com/golang-standards/project-layout/tree/master/website/README.md"><code>/website</code></a> pour des exemples</p>
<h2>Les répertoires que vous ne devriez pas avoir</h2>
<h3><code>/src</code></h3>
<p>Certains projets Go comportent un dossier <code>src</code> mais cela arrive en général lorsque les développeurs viennent du monde de Java où c'est une pratique habituelle. Faites tout votre possible pour ne pas adopter ce pattern Java. Vous n'avez vraiment pas envie que votre code Go ou vos projets Go ressemblent à du Java :-)</p>
<p>Ne confondez pas le répertoire <code>/src</code> à la racine avec le répertoire <code>/src</code> utilisé par Go pour gérer ses espaces de travail comme décrit dans <a href="https://golang.org/doc/code.html"><code>How to Write Go Code</code></a>. La variable d'environnement <code>$GOPATH</code> pointe vers votre espace de travail courant (par défault il pointe vers <code>$HOME/go</code> sur les systèmes non-Windows). Cet espace de travail inclut les dossiers <code>/pkg</code>, <code>/bin</code> et <code>/src</code>. Votre projet en lui-même va se retrouver dans un sous-dossier de <code>/src</code>, donc si vous avez un dossier <code>/src</code> dans votre projet le chemin vers celui-ci ressemblera à ceci : <code>/some/path/to/workspace/src/your_project/src/your_code.go</code>. Notez qu'à partir de Go 1.11 il est possible d'avoir votre projet en dehors de votre <code>GOPATH</code> mais cela ne veut toujours pas dire que c'est une bonne idée d'utiliser le dossier <code>/src</code></p>
<h2>Badges</h2>
<ul>
<li>
<p><a href="https://goreportcard.com/">Go Report Card</a> - Scanne votre code avec les commandes <code>gofmt</code>, <code>go vet</code>, <code>gocyclo</code>, <code>golint</code>, <code>ineffassign</code>, <code>license</code> and <code>misspell</code>. Remplacez <code>github.com/golang-standards/project-layout</code> avec l'url de votre projet.</p>
</li>
<li>
<p>~~<a href="http://godoc.org">GoDoc</a> - Fournit une version en ligne de votre documentation générée GoDoc. Modifiez le lien pour qu'il pointe vers votre projet.~~</p>
</li>
<li>
<p><a href="https://pkg.go.dev">Pkg.go.dev</a> - Pkg.go.dev est la nouvelle destination pour la découverte de Go et sa documentation. Vous pouvez créer une badge en utilisant <a href="https://pkg.go.dev/badge">l'outil de création de badge</a>.</p>
</li>
<li>
<p>Release - Il indique la dernière version de votre projet. Modifiez le lien GitHub pour qu'il pointe vers votre projet.</p>
</li>
</ul>
<p><a href="https://goreportcard.com/report/github.com/golang-standards/project-layout"><img src="https://goreportcard.com/badge/github.com/golang-standards/project-layout?style=flat-square" alt="Go Report Card"></a>
<a href="http://godoc.org/github.com/golang-standards/project-layout"><img src="https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square" alt="Go Doc"></a>
<a href="https://github.com/golang-standards/project-layout/releases/latest"><img src="https://img.shields.io/github/release/golang-standards/project-layout.svg?style=flat-square" alt="Release"></a></p>
<h2>Notes</h2>
<p>Un template de projet moins générique avec du code, des script et des configs réutilisables est en cours de réalisation.</p>

    </article>