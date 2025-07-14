How to run:
- For local without docker 
  - adjust key DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT, PORT, SECRET_KEY with local environment
  - just add permision to running.sh like : sudo chmod +x running.sh and run it
- For local with docker 
  - make sure you have installed docker compose
  - then enter command docker compose up --build
  - if have need permission from super user, just add sudo docker compose up --build
