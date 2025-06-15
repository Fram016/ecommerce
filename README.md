# Ecommerce-Go

Proyecto de Ecommerce desarrollado en Go, utilizando Gorilla Mux para el enrutamiento, sesiones seguras, y Bootstrap para el frontend.

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

El servidor estará disponible en [http://localhost:8080](http://localhost:8080).

## Funcionalidades

- Registro y login de usuarios (con roles: admin y cliente)
- CRUD de productos y categorías (solo admin)
- Gestión de direcciones de clientes
- Subida y gestión de imágenes de productos
- Vistas diferenciadas para admin y clientes

## Pendiente

- Gestión de perfil
- Carrito
- Panel de administración

## Recursos usados

- [Bootstrap](https://getbootstrap.com/) para estilos CSS
- [Gorilla Mux](https://github.com/gorilla/mux) para enrutamiento
- [Gorilla Sessions](https://github.com/gorilla/sessions) para manejo de sesiones
- [Go SQL Driver MySQL](https://github.com/go-sql-driver/mysql) para conexión a la base de datos

---

© 2025 Ecommerce-Go. Todos los derechos reservados.