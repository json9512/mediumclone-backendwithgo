# mediumclone-backendwithgo
[미디엄클론](https://github.com/json9512/mediumclone)의 백엔드를 Go로 재개발 하는 프로젝트입니다.

개발 과정과 의사 결정은 [블로그](https://json9512.github.io/blog/%ED%94%84%EB%A1%9C%EC%A0%9D%ED%8A%B8-Medium-%ED%81%B4%EB%A1%A0-%EB%B0%B1%EC%97%94%EB%93%9C%EB%A5%BC-%EB%A7%8C%EB%93%A4%EC%96%B4%EB%B3%B4%EC%9E%90-1/)에 기록중입니다.

백엔드를 Go lang으로 개발 중에 있습니다.

## 개발 목적

- Go lang 배우기
- RESTful한 API 만들어보기
- BDD/TDD 방식으로 개발해보기
- Docker, Github Actions, AWS 사용해보기

## 기술 스택

- [Go](https://golang.org/) 언어
    - [Gin](https://github.com/gin-gonic/gin) 웹 프레임워크
    - [Goblin](https://github.com/franela/goblin) 테스트 프레임워크
- AWS [RDS](https://aws.amazon.com/ko/rds/) (Postgresql)
- [Docker](https://www.docker.com/), Docker Hub
- Github Actions (CI/CD)
- AWS [ECR](https://aws.amazon.com/ko/ecr/) (Docker 이미지 배포 레지스트리)
- AWS [ECS](https://aws.amazon.com/ko/ecs/) (인스턴스 및 컨테이너 관리 서비스)
- AWS [EC2](https://aws.amazon.com/ko/ec2/) (서버가 배포되는 컴퓨팅 엔진)

