# wpscanサービス用のbase imageのビルドについて
CodeBuildのマネージドイメージに含まれているdockerのバージョンではwpscanのDockerfileのベースイメージをビルドできない。
そのため、以下のワークアラウンドで対応すること。
* amd64のEC2インスタンスでイメージをビルドしてECRにpush
* arm64のEC2インスタンスで同上
* pushしたイメージを元にmanifestを作成し、ECRにpush