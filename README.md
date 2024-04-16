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