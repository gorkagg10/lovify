# Lovify

## Preparación 

Para poder hacer uso de la aplicación es necesario tener instaladas las herramientas *Docker*,
*Docker Compose*, y *gpg*.

Para hacer uso de la aplicación se deberán seguir los siguientes pasos:

1. Clonar el código fuente de la aplicación ejecutando el siguiente comando:

  ```
  git clone https://github.com/gorkagg10/lovify.git
  ```

2. Acceder al directorio del código fuente utilizando el siguiente comando:

  ```
  cd lovify
   ```

3. Desencriptar cada uno de los volúmenes utilizados por docker-compose utilizando los
   siguientes comandos. La contraseña para todos los volúmenes será lovify:

  ```
  gpg --output mongodb-data.tar.gz --decrypt mongodb-data.tar.gz.gpg 
  gpg --output postgres-data.tar.gz --decrypt postgres-data.tar.gz.gpg 
  gpg --output nats-data.tar.gz --decrypt nats-data.tar.gz.gpg
  ```

4. Descomprimir cada uno de los volúmenes utilizados por docker-compose utilizando los siguientes comandos:

  ```
  tar -xzf mongodb-data.tar.gz
  tar -xzf postgres-data.tar.gz
  tar -xzf nats-data.tar.gz
   ```

5. Crear las imágenes de Docker con el siguiente comando:

  ```
  docker compose build
  ```

6. Poner en marcha la aplicación utilizando el comando:

  ```
  docker compose up -d
   ```

7. Acceder a la aplicación introduciendo la siguiente dirección en el navegador:
 http://localhost:3000

## Uso

Para la prueba de esta aplicación se disponen de dos usuarios:

- **admin**: con una cuenta vinculada de Spotify.
- **admin2**: al cual se le puede vincular una cuenta personal de Spotify.

Los datos del usuario **admin** son los siguientes:

- *Email:* admin@example.com
- *Contraseña:* soyeladmin10

Los datos del usuario admin2 son los siguientes:

- *Email:* admin2@example.com
- *Contraseña:* soyeladmin20 