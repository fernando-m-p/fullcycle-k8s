# Curso Fullcycle - Módulo Kubernetes

Esse projeto tem o objetiuvo de guardar códigos e anotações do conhecimento adiquirido no curso.

## 1. Andamento do curso

- [x] Docker
- [x] Padrões e técnicas avançadas com Git e Github
- [x] Integração contínua
- [ ] Kubernetes
- [ ] API Gateway
- [ ] API Gateway com Kong e Kubernetes
- [ ] Observabilidade
- [ ] Introdução a OpenTelemetry
- [ ] Terraform
- [ ] Ansible
- [ ] GitOps
- [ ] Deploy nas Cloud Providers
- [ ] Fundamentos da arquitetura de software
- [ ] Comunicação entre sistemas
- [ ] SOLID Express
- [ ] Domain Driven Design
- [ ] DDD: Modelagem Tática e Patterns
- [ ] Event Storming na Prática
- [ ] Arquitetura hexagonal
- [ ] Clean Architecture
- [ ] Sistemas monolíticos
- [ ] Arquitetura baseada em microsserviços
- [ ] EDA - Event Driven Architecture
- [ ] RabbitMQ
- [ ] Apache Kafka
- [ ] Autenticação e Keycloak
- [ ] Arquitetura do projeto prático - Codeflix
- [ ] Projeto prático - Java (Back-end)
- [ ] Projeto prático CodeFlix - React (Front-end)
- [ ] Projeto prático Admin CodeFlix - React (Front-end)
- [ ] Microsserviços de Encoder de Vídeo com Go Lang
- [ ] Microsserviços: API do Catálogo com Java (Back-end)
- [ ] MIcrosserviços: Assinaturas com Java (Back-end)

## 2. Anotações do Curso
### 2.1. Iniciando com kuberbetes


> Kubernetes (K8s) é um produto Open Source utilizado para automatizar a implantação, o dimensionamento e o gerenciamento de aplicativos em contêiner
>
> Fonte: [kubetentes.io](https://kubernetes.io/pt-br/)

No kubernetes utilizamos a API usando a CLI: **kubectl**


#### 2.1.1. Conceitos:
- Cluster: conjunto de máquinas (Nodes)
- Pods: Unidade que contém os containers provisionados.
- - O Pod representa os processos rodando no cluster, geralmente um pod roda somente um container, apesar de não ser regra.
- Deployment: Tem o objetivo de provisionar os Pods por meio dos ReplicaSets
- ReplicaSet: Tem como objetivo manter um conjunto estável de Pods de réplica em execução a quaquer momento. (usado para garantir a disponibilidade)

#### 2.1.2. Criando o primeiro cluster com o kind

O arquivo k8s/kind.yaml é um exemplo de criar um Cluster simples com 1 control-plane e 3 workers

é executado `kind create cluster --name <nome do cluster`

depois de rodar, para o kubectl ter acesso ao cluster devemos rodar
``` shell
 kubectl cluster-info --context <nome do cluster>
```

### 2.2. Primeiros passos na prática
#### 2.2.1. Aplicação de exemplo

Os arquivos `server.go` e `Dockerfile` são usados para criar a imagem da aplicação simples em Go que escreve uma frase quando é acessado o server.

#### 2.2.2. Trabalhando com Pods

O arquivo `k8s/pod.yaml` foi utilizado para a criação do primeiro pod, dessa vez como o cluster já está criado foi executado a instrução 

``` shell
kubectl apply -f k8s/pod.yaml
```


| atributo | descrição    |
| :-----: | :---:  |
| metadata.name | Nome do pod|
| metadata.labels | Esse campo é utilizado para filtrar pods para aplicar algumas configurações depois|
| spec.containers | Especificações dos containers que vai rodar dentro de cada pod|
| spec.containers.image | Imagem do container|

Foi adiantando um comando para acessar esse pod, fazendo um mapeamento da porta:

```shell
kubectl port-forward pod/{nome-pod} {LocalPort}:{PortPod}
```

#### 2.2.3. Criando primeira ReplicaSet

Os ReplicaSet tem como objetivo manter um conjunto estável de Pods de réplica em execução a quaquer momento. (usado para garantir a disponibilidade)

Com o arquivo `k8s/replicaset.yaml` foi criado o primeiro replicaset, colocando a especificação do pod anterior no atributo *spec.template* no atributo *spec.replicas* foi configurado quantos pods o replicaset deve garantir, caso um seja deletado ele recria outro com as mesmas especificações.

#### 2.2.4. Problema do ReplicaSet

O ReplicaSet garante a disponibilidade, mas caso uma alteração no arquivo do replicaSet (`k8s/replicaset.yaml`) seja feita e aplicada (`kubectl apply -f k8s/replicaset.yaml`) a alteração não é replicado nos pods, para a alteração ter efeito os pods devem ser deletados, assim quando for criado o novo pod terá a nova alteração.

#### 2.2.5. Implementando Deployment

Um objeto Deployment cria replicasets que gerenciam os pods, assim alterando o atributo kind para Deployment resolvemos o problema de uma alteração na configuração que não é replicada nos pods.

Depois de uma alteração no arquivo o kubernetes vai matando os pods do replicaset anterior e vai startando os novos pods do novo replicaset. Assim sem downtime.

Importante, o kubectl não mata o replicaset anterior apenas desativa os pods dele. Deixando ele disponível.

#### 2.2.6. Rollout e Revisões

Caso aconteça um erro em uma atualização do deployment e seja necessário retornar a versão anterior.

**Para ver o histórico de versões**

```shell
kubectl rollout history deployment goserver
```

Para escrever uma mensagem no CHANGE-CAUSE adicione o atributo `metadata.annotations.kubernetes.io/change-cause` no seu arquivo .yaml

**Para voltar para versão anterior**

```shell
    kubectl rollout undo deployment goserver
```

**Para voltar para uma versão específica x**


```shell
    kubectl rollout undo deployment goserver --to-revision=x
```

### 2.3. Services

#### 2.3.1. Conceito

Service é um método para expor um aplicativo de rede que está sendo executado como um ou mais pods no seu cluster.

No curso é apresentado os seguintes tipos de services:
- ClusterIP
- NodePort
- LoadBalancer

##### ClusterIP

Expõe o Serviço em um IP interno do cluster. Escolhendo este valor torna o Serviço acessível apenas a partir do cluster. Este é o padrão que é usado se você não especificar explicitamente um para um Serviço. Você pode expor o Serviço à Internet pública usando Ingress ou um Gateway.type

Exemplo:

``` yaml
    apiVersion: v1
    kind: Service
    metadata:
    name: goserver-service
    spec:
    selector:
        app: goserver
    type: ClusterIP
    ports:
    - name: goserver-service
        port: 80
        targetPort: 8000
        protocol: TCP
```

Lembrando que o atributo `port` é a porta de entrada do Service e `portTarget` é a saída do Service e entrada dos pods.


##### NodePort

O uso de um NodePort lhe dá a liberdade de configurar sua própria solução de balanceamento de carga, para configurar ambientes que não são totalmente suportados pelo Kubernetes, ou mesmo para expor os endereços IP de um ou mais nós diretamente.

Cada nó no cluster configura para ouvir na porta atribuída e encaminhar o tráfego para um dos prontos pontos de extremidade associados a esse Serviço.

obs: a porta especificada no atributo spec.ports.nodePort por padrão deve ser especificado no intervalo (30000-32767)

##### LoadBalancer

Em provedores de nuvem que oferecem suporte a balanceadores de carga externos, defina o campo para provisionar um balanceador de carga para seu Serviço. A criação real do balanceador de carga acontece de forma assíncrona, e as informações sobre o balanceador provisionado são publicadas no campo Serviço.

O tráfego do balanceador de carga externo é direcionado para os Pods de back-end. A nuvem O provedor decide como ele é balanceado de carga.


### 2.4. Objetos de configuração

#### 2.4.1. Variáveis de Ambiente

A primeira forma de trabalhar com variáveis de ambiente é adicionar a variável de ambiente no arquivo .yaml do Deployment

Especificando no atributo `spec.template.spec.containers.env`, mas não é uma boa prática utilizar

Versão do Deployment utilizando as variáveis de ambientes HardCode
```yaml

apiVersion: apps/v1
kind: Deployment

metadata:
  name: goserver
  labels:
    apps: goserver
  annotations:
    kubernetes.io/change-cause: "Deployment: goserver com imagem v1"

spec:
  selector:
    matchLabels:
      app: goserver
  replicas: 2
  template:
    metadata:
      labels:
        app: "goserver"
    spec:
      containers:
      - name: goserver
        image: "fermope/go-server:v3"
        env:
        - name: NAME
          value: "Fernando"
        - name: AGE
          value: "33"
```

#### 2.4.2. Utilizando ConfigMap


Um ConfigMap é um objeto de API usado para armazenar dados não confidenciais em pares chave-valor. Pod pode consumir ConfigMaps como variáveis de ambiente, argumentos de linha de comando ou como arquivos de configuração em um volume.

Há quatro maneiras diferentes de usar um ConfigMap para configurar um recipiente dentro de um Pod:

1. Dentro de um contêiner comando e args
2. Variáveis de ambiente para um contêiner
3. Adicionar um arquivo no volume somente leitura, para o aplicativo ler
4. Escrever código para ser executado dentro do Pod que usa a API do Kubernetes para ler um ConfigMap

Para os três primeiros métodos, o Kubelet usa os dados de o ConfigMap quando ele inicia contêiner(es) para um Pod.

``` yaml
apiVersion: v1
kind: Pod
metadata:
  name: configmap-demo-pod
spec:
  containers:
    - name: demo
      image: alpine
      command: ["sleep", "3600"]
      env:
        # Define the environment variable
        - name: PLAYER_INITIAL_LIVES # Notice that the case is different here
                                     # from the key name in the ConfigMap.
          valueFrom:
            configMapKeyRef:
              name: game-demo           # The ConfigMap this value comes from.
              key: player_initial_lives # The key to fetch.
        - name: UI_PROPERTIES_FILE_NAME
          valueFrom:
            configMapKeyRef:
              name: game-demo
              key: ui_properties_file_name
      volumeMounts:
      - name: config
        mountPath: "/config"
        readOnly: true
  volumes:
  # You set volumes at the Pod level, then mount them into containers inside that Pod
  - name: config
    configMap:
      # Provide the name of the ConfigMap you want to mount.
      name: game-demo
      # An array of keys from the ConfigMap to create as files
      items:
      - key: "game.properties"
        path: "game.properties"
      - key: "user-interface.properties"
        path: "user-interface.properties"
        
```

#### 2.4.3. Secrets e variáveis de ambiente

Um Segredo é um objeto que contém uma pequena quantidade de dados confidenciais, como uma senha, um token ou uma chave. Tais informações poderiam, de outra forma, ser colocadas em um Vagem especificação ou em um imagem do contêiner. Usando um Segredo significa que você não precisa incluir dados confidenciais em seu código do aplicativo.

Como os Segredos podem ser criados independentemente dos Pods que os usam, há é menor o risco de o Segredo (e seus dados) serem expostos durante o fluxo de trabalho de criação, visualização e edição de Pods. Kubernetes e aplicativos que são executados em seu cluster, também pode tomar precauções adicionais com Secrets, como evitar gravar dados confidenciais em armazenamento não volátil.

o arquivo `k8s/secret.yaml` é um exemplo de um Secret do tipo Opaque

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: goserver-secret
type: Opaque
data:
  USER: "RmVybmFuZG8K" #Os segredos são adicionados em Base64
  PASSWORD: "MTIzNDU2Cg=="
```

[Veja mais](https://kubernetes.io/docs/concepts/configuration/secret/)


### 2.5. Probes

Em vários cenários é muito importante termos uma forma de informar se a aplicação está rodando, seja para avisar a equipe de manutenção ou seja para realizar tarefas automatizadas.

#### 2.5.1. LivenessProbe

Com o LivenessProbe podemos reiniciar os Pods da nossa aplicação, baseado em uma verificação periódica.

```yaml
[...]
spec:
  template:
    spec:
        containers:
        - name: goserver
          image: "fermope/go-server:v5"
          livenessProbe:
            httpGet:
              path: /healthz
              port: 8000
            periodSeconds: 5 # Periódicidade da verificação
            failureThreshold: 3 # Numero de verificações para ser considerada uma falha na aplicação e ser reiniciado o Pod
            timeoutSeconds: 1 # Timeout para a verificação
            successThreshold: 1 # Numero de verificações para a aplicação ser considerada saudável
[...]
```

#### 2.5.2. ReadinessProbe

Com o ReadinessProbe podemos podemos verificar se a nossa aplicação está saudável antes de encaminhar o tráfego para os Pods.

```yaml
[...]
spec:
  template:
    spec:
        containers:
        - name: goserver
          readinessProbe:
            httpGet:
              path: /healthz
              port: 8000
            periodSeconds: 3 # Periódicidade da verificação
            failureThreshold: 1 # Numero de verificações para ser considerada uma falha na aplicação
            initialDelaySeconds: 10 # Quanto tempo inicial para iniciar as verificações
[...]
```

#### 2.5.3. StartupProbe

O StartupProbe é utilizado para verificar se a aplicação iniciou colocando para verificar periódicamente e com um FailThreshoud alto é uma boa técnica para substituir o initialDelaySeconds.

Ou seja com o initialDelaySeconds de 90 segundos a aplicação só vai ser testada depois dos 90 segundos, com o startupProbe com o periodSeconds 3 e o failThreshoud 30 vai ser verificado os 90 segundos e quando a aplicação ficar saudável vai livberar o tráfego, tornando assim a aplicação disponível entre 3 e no máximo 90 segundos.


```yaml
spec:
  template:
    spec:
        containers:
        - name: goserver
          startupProbe:
            httpGet:
              path: /healthz
              port: 8000
            periodSeconds: 3
            failureThreshold: 30
```