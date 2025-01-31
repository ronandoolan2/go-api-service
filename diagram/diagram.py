from diagrams import Diagram, Cluster, Edge
from diagrams.aws.general import User
from diagrams.k8s.compute import Pod
from diagrams.k8s.ecosystem import Helm
from diagrams.k8s.network import Service
from diagrams.onprem.monitoring import Prometheus
from diagrams.onprem.database import PostgreSQL

with Diagram("Kubernetes Transaction API Architecture", show=False):
    # Represent the external user
    user = User("User")

    # High-level cluster box
    with Cluster("Kubernetes Cluster"):
        # Prometheus Operator
        operator = Helm("Prometheus Operator")

        # Transaction API namespace / pods
        with Cluster("API Namespace"):
            api_service = Service("transaction-api-service")
            api_pod = Pod("transaction-api-pod")

        # Database namespace / pods
        with Cluster("DB Namespace"):
            postgres_service = Service("postgres-service")
            postgres_pod = Pod("postgres-pod")

        # Prometheus
        with Cluster("Monitoring Namespace"):
            prometheus = Prometheus("Prometheus")

        # Flow:
        # User -> Service -> Pod -> DB
        user >> Edge(label="makes request") >> api_service >> api_pod >> Edge(label="DB connection") >> postgres_service
        postgres_service >> postgres_pod

        # Prometheus scrapes metrics from the API and Postgres (via exporters)
        prometheus >> Edge(label="scrapes metrics") >> api_service
        prometheus >> Edge(label="scrapes metrics") >> postgres_service

        # Operator manages Prometheus
        operator >> prometheus
