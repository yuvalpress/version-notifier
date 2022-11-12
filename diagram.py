from diagrams import Cluster, Diagram
from diagrams.saas.chat import Slack
from diagrams.programming.language import Go
from diagrams.onprem.vcs import Github
from diagrams.custom import Custom

with Diagram("\nVersion Notifier", show=False):
    with Cluster(""):
        repo = Github("Repository")
        
        with Cluster("Version Notifier"):
            notifier = Go("Version Notifier")
            custom = Custom(label="Repo List", icon_path="./docs/config-file.png")
            custom - notifier
        
        notifier >> repo
        repo >> notifier
        notifier >> Slack("Slack")