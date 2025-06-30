# Ecommerce-Go

- **Estudiante:** Francisco López  
- **Curso:** Programación Orientada a Objetos
- **Fecha:** Junio 2024

## Objetivo del Programa

El objetivo de este proyecto es desarrollar una plataforma de comercio electrónico (Ecommerce) utilizando Go como lenguaje principal, que permita la gestión de productos, categorías, usuarios y direcciones, diferenciando funcionalidades entre administradores y clientes. El sistema busca ser seguro, modular y fácil de mantener, integrando buenas prácticas de desarrollo web y seguridad.

## Principales Funcionalidades

- **Registro y login de usuarios:**  
  Permite a los usuarios crear una cuenta y autenticarse en el sistema. Se manejan roles de usuario (admin y cliente) para controlar el acceso a diferentes funcionalidades.

- **Gestión de productos y categorías (CRUD):**  
  Los administradores pueden crear, leer, actualizar y eliminar productos y categorías. Los productos incluyen información como nombre, descripción, precio, stock, imágenes y categoría.

- **Gestión de direcciones de clientes:**  
  Los clientes pueden agregar, editar y eliminar sus direcciones de envío, eligiendo una dirección principal para sus pedidos.

- **Subida y gestión de imágenes de productos:**  
  Los administradores pueden subir imágenes para los productos, facilitando una mejor presentación visual en la tienda.

- **Vistas diferenciadas para admin y clientes:**  
  El sistema muestra diferentes interfaces y opciones según el rol del usuario, asegurando que solo los administradores accedan a la gestión interna.

- **Panel de administración:**  
  Acceso exclusivo para administradores, donde pueden gestionar usuarios, productos, categorías y pedidos.

- **Carrito de compras (pendiente):**  
  Funcionalidad planificada para permitir a los clientes agregar productos a un carrito y realizar pedidos.

## Requisitos

- Go 1.24.3 o superior
- MariaDB/MySQL
- Node.js y npm (para instalar dependencias de frontend)

## Instalación

1. Clona el repositorio y entra en la carpeta del proyecto.
2. Copia el archivo `.env` y configura tus credenciales de base de datos.
3. Instala las dependencias de Node.js para Bootstrap:

   ```sh
   npm install
   ```

4. Instala las dependencias de Go:

   ```sh
   go mod download
   ```

5. Asegúrate de tener la base de datos creada y migrada.

## Estructura del Proyecto

- `main.go`: Punto de entrada de la aplicación.
- `controllers/`: Lógica de controladores para usuarios, productos, categorías, login, etc.
- `models/`: Modelos de datos y acceso a la base de datos.
- `middlewares/`: Middlewares para autenticación y autorización.
- `db/`: Conexión a la base de datos.
- `views/`: Plantillas HTML para las vistas.
- `public/`: Archivos públicos (imágenes, etc).
- `static/`: Archivos estáticos servidos desde node_modules (Bootstrap).

## Uso

Para iniciar el servidor:

```sh
go run main.go
```

El servidor estará disponible en [http://localhost:8081](http://localhost:8081).

## Herramientas usadas

- **Driver mysql de GO**: Para conectar a la base de datos (compatible con MariaDB).
- **Gotgotenv**: Para gestionar credenciales y variables de entorno.
- **Mux**: Para la gestión de rutas HTTP.
- **NPM**: Para gestión de dependencias del frontend.
- **Bootstrap**: Para el diseño responsive de la interfaz.
- **Bcrypt**: Para la encriptación segura de contraseñas.
- **Gorilla Sessions**: Para manejo seguro de sesiones y cookies.
- **Templates**: Para el frontend modular con plantillas HTML en Go.

---

© 2025 Ecommerce-Go. Todos los derechos reservados.