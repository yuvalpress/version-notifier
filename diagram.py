from diagrams import Cluster, Diagram
from diagrams.saas.chat import Slack, Telegram
from diagrams.programming.language import Go
from diagrams.onprem.vcs import Github
from diagrams.custom import Custom

with Diagram("\nVersion Notifier Service", show=False, direction="LR"):
    with Cluster(""):
        repo = Github("Repositories")
        
        with Cluster("Version Notifier", direction="LR"):
            notifier = Go("Code")
            custom = Custom(label="Repo List", icon_path="./docs/config-file.png")
            custom - notifier
        
        notifier >> repo
        repo >> notifier

        with Cluster("Channels", direction="lr"):
            slack = Slack("Slack")
            telegram = Telegram("Telegram")

        notifier >> slack
