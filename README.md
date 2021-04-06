# mediumclone-backendwithgo
[미디엄클론](https://github.com/json9512/mediumclone)의 백엔드를 Go로 재개발 하는 프로젝트입니다.

개발 과정과 의사 결정은 [블로그](https://json9512.github.io/blog/%ED%94%84%EB%A1%9C%EC%A0%9D%ED%8A%B8-Medium-%ED%81%B4%EB%A1%A0-%EB%B0%B1%EC%97%94%EB%93%9C%EB%A5%BC-%EB%A7%8C%EB%93%A4%EC%96%B4%EB%B3%B4%EC%9E%90-1/)에 기록중입니다.

백엔드를 Go lang으로 개발 중에 있습니다.

## 개발 목적
- Golang 배우기
- RESTful한 API 만들어보기
- BDD/TDD 방식으로 개발해보기
- Docker, Github Actions, AWS 사용해보기

## 기술 스택
- [Go](https://golang.org/) 개발 언어
    - [Gin](https://github.com/gin-gonic/gin) 웹 프레임워크
    - [Goblin](https://github.com/franela/goblin) 테스트 프레임워크
- AWS [RDS](https://aws.amazon.com/ko/rds/) Postgresql
- [Docker](https://www.docker.com/), Docker Hub
- Github Actions (CI/CD)
- AWS [ECR](https://aws.amazon.com/ko/ecr/) Docker 이미지 배포 레지스트리
- AWS [ECS](https://aws.amazon.com/ko/ecs/) 인스턴스 및 컨테이너 관리 서비스
- AWS [EC2](https://aws.amazon.com/ko/ec2/) 서버가 배포되는 컴퓨팅 엔진

## Endpoints
### `/posts`
<br>
<details close>

<summary> </summary>

####  GET
모든 `posts`를 리턴함

#### POST
새로운 `post`를 생성함.
<details close>

<summary>Request Body</summary>

- `doc` (required, JSON): `post` 본문
- `tags` (string): 연관 tag를 `,` 로 구분. ex `tags: sports, soccer, football`
- `likes` (integer): 좋아요 수
- `comments` (JSON): 관련 comments.
</details>
<br>

#### PUT
기존 `post`를 업데이트함. `id`외에도 `doc`, `tags`, `likes`, `comments` 중 하나 이상이 필요함
<details close>

<summary>Request Body</summary>

- `id` (required, int): `post`의 ID
- `doc` (JSON): `post` 본문
- `tags` (string): 연관 tag를 `,` 로 구분. ex `tags: sports, soccer, football`
- `likes` (integer): 좋아요 수
- `comments` (JSON): 관련 comments.
</details>
<br>

</details>
<br>

### `/posts/:id`
<br>
<details>

<summary> </summary>

#### GET
주어진 `id`의 `post`를 리턴함

#### DELETE
주어진 `id`의 `post`를 삭제함

</details>
<br>

### `/posts/:id/like`
<br>
<details>

<summary> </summary>

#### GET
주어진 `id`의 `post`의 `like`를 리턴함

</details>
<br>

### `/users`
<br>
<details>

<summary> </summary>

#### POST
새로운 `user`를 생성함
<details close>

<summary>Request Body</summary>

- `email` (required, string): email
- `password` (required, string): password
</details>
<br>

#### PUT
기존 `user` 정보를 수정함. `id` 외에도 `email` 이나 `password`가 필요함
<details close>

<summary>Request Body</summary>

- `id` (required, int): `user`의 ID
- `email` (string): email
- `password` (string): password
</details>
<br>

</details>
<br>

### `/users/:id`
<br>
<details>

<summary> </summary>

#### GET
주어진 `id`의 `user`를 리턴함

#### DELETE
주어진 `id`의 `users`를 삭제함

</details>
<br>

### `/login`
<br>
<details>

<summary> </summary>

#### POST
주어진 정보로 `user`를 인증함

`httpOnly` Cookie에 `access_token`을 저장함
<details close>

<summary>Request Body</summary>

- `email` (required, string): email
- `password` (required, string): password
</details>
<br>
</details>
<br>

### `/logout`
<br>
<details>

<summary> </summary>

#### POST
주어진 `email`로 `user`를 로그아웃함

`httpOnly` Cookie에 `access_token`을 초기화함
<details close>

<summary>Request Body</summary>

- `email` (required, string): email
</details>
<br>
</details>
<br>