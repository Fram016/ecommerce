/*
 Navicat Premium Data Transfer

 Source Server         : local
 Source Server Type    : MariaDB
 Source Server Version : 110404 (11.4.4-MariaDB)
 Source Host           : localhost:5588
 Source Schema         : ecommerce

 Target Server Type    : MariaDB
 Target Server Version : 110404 (11.4.4-MariaDB)
 File Encoding         : 65001

 Date: 15/06/2025 13:09:14
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for categorias
-- ----------------------------
DROP TABLE IF EXISTS `categorias`;
CREATE TABLE `categorias`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `nombre` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `descripcion` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_nombre_categoria`(`nombre` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for clientes_direcciones
-- ----------------------------
DROP TABLE IF EXISTS `clientes_direcciones`;
CREATE TABLE `clientes_direcciones`  (
  `id` int(20) NOT NULL AUTO_INCREMENT,
  `usuario_id` int(11) NULL DEFAULT NULL,
  `direccion` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `ciudad` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `provincia` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `codigo_postal` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `pais` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `es_principal` tinyint(1) NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_usuario_id`(`usuario_id` ASC) USING BTREE,
  CONSTRAINT `fk_usuario_direccion` FOREIGN KEY (`usuario_id`) REFERENCES `usuarios` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for pedidos
-- ----------------------------
DROP TABLE IF EXISTS `pedidos`;
CREATE TABLE `pedidos`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `usuario_id` int(11) NULL DEFAULT NULL,
  `fecha_pedido` timestamp NULL DEFAULT current_timestamp(),
  `estado` enum('pendiente','procesado','enviado','entregado','cancelado') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'pendiente',
  `total` decimal(10, 2) NOT NULL,
  `direccion_envio` int(20) NOT NULL,
  `observacion` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `fk_usuario`(`usuario_id` ASC) USING BTREE,
  INDEX `idx_estado_pedido`(`estado` ASC) USING BTREE,
  INDEX `fk_clientes_direcciones`(`direccion_envio` ASC) USING BTREE,
  CONSTRAINT `fk_usuario` FOREIGN KEY (`usuario_id`) REFERENCES `usuarios` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_clientes_direcciones` FOREIGN KEY (`direccion_envio`) REFERENCES `clientes_direcciones` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for pedidos_detalles
-- ----------------------------
DROP TABLE IF EXISTS `pedidos_detalles`;
CREATE TABLE `pedidos_detalles`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pedido_id` int(11) NULL DEFAULT NULL,
  `producto_id` int(11) NULL DEFAULT NULL,
  `cantidad` int(11) NOT NULL,
  `precio_unitario` decimal(10, 2) NOT NULL,
  `observacion` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `fk_pedido`(`pedido_id` ASC) USING BTREE,
  INDEX `fk_producto_detalle`(`producto_id` ASC) USING BTREE,
  CONSTRAINT `fk_pedido` FOREIGN KEY (`pedido_id`) REFERENCES `pedidos` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_producto_detalle` FOREIGN KEY (`producto_id`) REFERENCES `productos` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for producto_imagenes
-- ----------------------------
DROP TABLE IF EXISTS `producto_imagenes`;
CREATE TABLE `producto_imagenes`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `producto_id` int(11) NULL DEFAULT NULL,
  `ruta_imagen` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `tipo_imagen` enum('principal','galeria') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'galeria',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_producto_id`(`producto_id` ASC) USING BTREE,
  CONSTRAINT `fk_producto` FOREIGN KEY (`producto_id`) REFERENCES `productos` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 25 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for productos
-- ----------------------------
DROP TABLE IF EXISTS `productos`;
CREATE TABLE `productos`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `nombre` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `descripcion` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `precio` decimal(10, 2) NOT NULL,
  `stock` int(11) NULL DEFAULT 0,
  `categoria_id` int(11) NULL DEFAULT NULL,
  `fecha_agregado` timestamp NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `fk_categoria`(`categoria_id` ASC) USING BTREE,
  INDEX `idx_nombre`(`nombre` ASC) USING BTREE,
  CONSTRAINT `fk_categoria` FOREIGN KEY (`categoria_id`) REFERENCES `categorias` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for usuarios
-- ----------------------------
DROP TABLE IF EXISTS `usuarios`;
CREATE TABLE `usuarios`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `correo` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `nombres` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `clave_segura` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `rol` enum('cliente','admin') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'cliente',
  `fecha_registro` timestamp NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `correo`(`correo` ASC) USING BTREE,
  CONSTRAINT `chk_correo` CHECK (`correo` like '%@%')
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
