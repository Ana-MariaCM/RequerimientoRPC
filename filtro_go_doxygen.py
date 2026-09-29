#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
@file filtro_go_doxygen.py
@brief Filtro de entrada (INPUT_FILTER) que permite documentar código Go con Doxygen.

Doxygen no reconoce la sintaxis de Go. Este filtro convierte cada archivo .go en
declaraciones con sintaxis tipo C++ que Doxygen sí interpreta, conservando los
comentarios de documentación (/ ** ... * /, ///<) y el número de cada línea:

  - package x                -> namespace Modulo { namespace x {
  - type T struct { ... }    -> struct T : public Embebido { Tipo campo; ... };
  - type I interface { ... } -> class I { public: virtual R Metodo(...) = 0; };
  - func (r *T) M(p) R {...} -> R M(p);   (con @memberof T)
  - func F(p) R {...}        -> R F(p);
  - const / var              -> declaraciones de variables

Uso (lo invoca Doxygen): python3 filtro_go_doxygen.py archivo.go
"""
import os
import re
import sys

IDENT = r"[A-Za-z_][A-Za-z0-9_]*"


def dividir_nivel_superior(texto, separador=","):
    """Divide 'texto' por 'separador' ignorando los que están dentro de (), [] o {}."""
    partes, nivel, actual = [], 0, ""
    for c in texto:
        if c in "([{":
            nivel += 1
        elif c in ")]}":
            nivel -= 1
        if c == separador and nivel == 0:
            partes.append(actual)
            actual = ""
        else:
            actual += c
    if actual.strip():
        partes.append(actual)
    return [p.strip() for p in partes if p.strip()]


def convertir_tipo(tipo):
    """Traduce un tipo de Go a una notación legible para Doxygen."""
    tipo = tipo.strip()
    if not tipo:
        return ""
    if tipo.startswith("..."):
        return "slice<%s>" % convertir_tipo(tipo[3:])
    if tipo.startswith("*"):
        return convertir_tipo(tipo[1:]) + "*"
    if tipo.startswith("[]"):
        return "slice<%s>" % convertir_tipo(tipo[2:])
    if tipo.startswith("map["):
        nivel, i = 0, 3
        while i < len(tipo):
            if tipo[i] == "[":
                nivel += 1
            elif tipo[i] == "]":
                nivel -= 1
                if nivel == 0:
                    break
            i += 1
        return "map<%s, %s>" % (convertir_tipo(tipo[4:i]), convertir_tipo(tipo[i + 1:]))
    for prefijo in ("<-chan ", "chan<- ", "chan "):
        if tipo.startswith(prefijo):
            return "chan<%s>" % convertir_tipo(tipo[len(prefijo):])
    if tipo.startswith("func"):
        return "func"
    if tipo in ("interface{}", "any"):
        return "any"
    if tipo == "struct{}":
        return "empty_struct"
    return tipo.replace(".", "::")


def convertir_parametros(texto):
    """Convierte 'a, b int, c string' en 'int a, int b, string c'."""
    elementos = dividir_nivel_superior(texto)
    resultado, pendientes = [], []
    for elemento in elementos:
        partes = elemento.split(None, 1)
        if len(partes) == 2 and re.fullmatch(IDENT, partes[0]):
            tipo = convertir_tipo(partes[1])
            for nombre in pendientes:
                resultado.append("%s %s" % (tipo, nombre))
            pendientes = []
            resultado.append("%s %s" % (tipo, partes[0]))
        else:
            pendientes.append(elemento)
    for tipo in pendientes:  # parámetros sin nombre
        resultado.append(convertir_tipo(tipo))
    return ", ".join(resultado)


def convertir_resultados(texto):
    """Convierte los resultados de una función de Go en un tipo de retorno."""
    texto = texto.strip()
    if not texto:
        return "void"
    if texto.startswith("(") and texto.endswith(")"):
        elementos = dividir_nivel_superior(texto[1:-1])
        tipos = []
        for elemento in elementos:
            partes = elemento.split(None, 1)
            tipos.append(convertir_tipo(partes[1] if len(partes) == 2 and re.fullmatch(IDENT, partes[0]) else elemento))
        return tipos[0] if len(tipos) == 1 else "tuple<%s>" % ", ".join(tipos)
    return convertir_tipo(texto)


def separar_parentesis(texto, inicio):
    """Devuelve (contenido, índice_final) del paréntesis que abre en 'inicio'."""
    nivel = 0
    for i in range(inicio, len(texto)):
        if texto[i] == "(":
            nivel += 1
        elif texto[i] == ")":
            nivel -= 1
            if nivel == 0:
                return texto[inicio + 1:i], i
    return texto[inicio + 1:], len(texto)


def quitar_cadenas_y_comentarios(linea):
    """Elimina literales y comentarios de una línea para contar llaves."""
    resultado, i, n = "", 0, len(linea)
    while i < n:
        c = linea[i]
        if linea.startswith("//", i):
            break
        if c in "\"'`":
            fin = i + 1
            while fin < n and linea[fin] != c:
                fin += 2 if linea[fin] == "\\" and c != "`" else 1
            i = fin + 1
            continue
        resultado += c
        i += 1
    return resultado


def posicion_comentario(texto):
    """Devuelve la posición de '//' fuera de literales de cadena, o -1."""
    i, n = 0, len(texto)
    while i < n:
        c = texto[i]
        if texto.startswith("//", i):
            return i
        if c in "\"'`":
            fin = i + 1
            while fin < n and texto[fin] != c:
                fin += 2 if texto[fin] == "\\" and c != "`" else 1
            i = fin + 1
            continue
        i += 1
    return -1


def inferir_tipo(valor):
    valor = valor.strip()
    if valor.startswith('"') or valor.startswith("`"):
        return "string"
    if re.match(r"^-?\d", valor):
        return "float" if re.match(r"^-?\d+\.\d", valor) else "int"
    if valor in ("true", "false"):
        return "bool"
    if valor.startswith("errors.New"):
        return "error"
    return "auto"


def declaracion_variable(texto, es_constante):
    """Convierte 'nombre [Tipo] [= valor] [comentario]' en una declaración."""
    comentario = ""
    posicion = posicion_comentario(texto)
    if posicion >= 0:
        comentario = " " + texto[posicion:]
        texto = texto[:posicion]
    texto = texto.strip()
    if not texto:
        return comentario.strip()
    valor = ""
    if "=" in texto:
        texto, valor = [t.strip() for t in texto.split("=", 1)]
    partes = texto.split(None, 1)
    nombre = partes[0]
    tipo = convertir_tipo(partes[1]) if len(partes) == 2 else inferir_tipo(valor)
    prefijo = "const " if es_constante else ""
    if valor and not valor.rstrip().endswith("{"):
        return "%s%s %s = %s;%s" % (prefijo, tipo, nombre, valor, comentario)
    return "%s%s %s;%s" % (prefijo, tipo, nombre, comentario)


def filtrar(ruta):
    with open(ruta, encoding="utf-8") as archivo:
        lineas = archivo.read().split("\n")

    raiz = os.path.dirname(os.path.abspath(__file__))
    relativa = os.path.relpath(os.path.abspath(ruta), raiz).split(os.sep)
    modulo = re.sub(r"\W", "_", relativa[0]) if len(relativa) > 1 else "raiz"

    salida = []
    estado = None            # None, 'comentario', 'import', 'struct', 'interface', 'bloque_var', 'cuerpo', 'cabecera'
    es_constante = False
    profundidad = 0
    indice_struct = -1
    cabecera = []
    ultimo_doc = -1          # índice de la última línea que cierra un comentario /** */
    paquete_abierto = False

    def emitir_funcion(texto):
        """Traduce la cabecera completa de una función y devuelve la declaración."""
        texto = texto.strip()
        texto = texto[len("func"):].strip()
        receptor = None
        if texto.startswith("("):
            contenido, fin = separar_parentesis(texto, 0)
            partes = contenido.split()
            receptor = partes[-1].lstrip("*") if partes else None
            receptor = re.sub(r"\[.*\]", "", receptor or "")
            texto = texto[fin + 1:].strip()
        m = re.match(r"(%s)\s*(\[[^\]]*\])?\s*\(" % IDENT, texto)
        if not m:
            return ""
        nombre = m.group(1)
        parametros, fin = separar_parentesis(texto, m.end() - 1)
        resultados = texto[fin + 1:].strip()
        if resultados.endswith("{"):
            resultados = resultados[:-1]
        declaracion = "%s %s(%s);" % (convertir_resultados(resultados), nombre, convertir_parametros(parametros))
        if receptor:
            if ultimo_doc >= 0 and all(not l.strip() for l in salida[ultimo_doc + 1:]):
                salida[ultimo_doc] = salida[ultimo_doc].replace("*/", "@memberof %s */" % receptor, 1)
            else:
                declaracion = "/** @memberof %s */ %s" % (receptor, declaracion)
        return declaracion

    for linea in lineas:
        limpia = linea.strip()

        # ----- comentarios de bloque (se conservan) -----
        if estado == "comentario":
            salida.append(linea)
            if "*/" in linea:
                estado = estado_previo
                if es_doc:
                    ultimo_doc = len(salida) - 1
            continue
        if limpia.startswith("/*") and estado not in ("cuerpo", "cabecera"):
            salida.append(linea)
            es_doc = limpia.startswith("/**") or limpia.startswith("/*!")
            if "*/" in limpia[2:]:
                if es_doc:
                    ultimo_doc = len(salida) - 1
            else:
                estado_previo, estado = estado, "comentario"
            continue

        # ----- cuerpo de funciones (se omite) -----
        if estado == "cuerpo":
            profundidad += quitar_cadenas_y_comentarios(linea).count("{") - quitar_cadenas_y_comentarios(linea).count("}")
            salida.append("")
            if profundidad <= 0:
                estado = None
            continue

        # ----- cabecera de función en varias líneas -----
        if estado == "cabecera":
            cabecera.append(limpia)
            salida.append("")
            texto = " ".join(cabecera)
            sin = quitar_cadenas_y_comentarios(texto)
            if sin.count("(") == sin.count(")") and sin.rstrip().endswith("{"):
                salida[indice_struct] = emitir_funcion(texto)
                profundidad = sin.count("{") - sin.count("}")
                estado = "cuerpo" if profundidad > 0 else None
            continue

        # ----- import -----
        if estado == "import":
            salida.append("")
            if limpia.startswith(")"):
                estado = None
            continue
        if limpia.startswith("import"):
            salida.append("")
            if limpia.endswith("("):
                estado = "import"
            continue

        # ----- package -----
        m = re.match(r"package\s+(%s)" % IDENT, limpia)
        if m and estado is None:
            salida.append("namespace %s { namespace %s {" % (modulo, m.group(1)))
            paquete_abierto = True
            continue

        # ----- campos de struct -----
        if estado == "struct":
            if limpia.startswith("}"):
                salida.append("};")
                estado = None
                continue
            texto = re.sub(r"`[^`]*`", "", linea)
            comentario = ""
            posicion = posicion_comentario(texto)
            if posicion >= 0:
                comentario = texto[posicion:].strip()
                texto = texto[:posicion]
            partes = texto.split()
            if not partes:
                salida.append(comentario)
            elif len(partes) == 1:  # campo embebido -> herencia
                base = convertir_tipo(partes[0]).rstrip("*")
                separador = " : public " if ":" not in salida[indice_struct].split("{")[0].replace("::", "") else ", public "
                salida[indice_struct] = salida[indice_struct].replace(" {", "%s%s {" % (separador, base), 1)
                salida.append(comentario)
            else:
                nombres = [n.strip() for n in " ".join(partes).split(",")]
                ultimo = nombres[-1].split(None, 1)
                nombres[-1] = ultimo[0]
                tipo = convertir_tipo(ultimo[1]) if len(ultimo) > 1 else ""
                salida.append("%s %s; %s" % (tipo, ", ".join(nombres), comentario))
            continue

        # ----- métodos de interfaz -----
        if estado == "interface":
            if limpia.startswith("}"):
                salida.append("};")
                estado = None
                continue
            texto = limpia.split("//")[0].strip()
            m = re.match(r"(%s)\s*\(" % IDENT, texto)
            if m:
                parametros, fin = separar_parentesis(texto, m.end() - 1)
                salida.append("virtual %s %s(%s) = 0;" % (convertir_resultados(texto[fin + 1:]), m.group(1),
                                                          convertir_parametros(parametros)))
            elif texto:
                base = convertir_tipo(texto)
                salida[indice_struct] = salida[indice_struct].replace(" {", " : public %s {" % base, 1)
                salida.append("")
            else:
                salida.append(linea)
            continue

        # ----- bloques const ( ... ) / var ( ... ) -----
        if estado == "bloque_var":
            if limpia.startswith(")"):
                salida.append("")
                estado = None
            else:
                salida.append(declaracion_variable(limpia, es_constante) if limpia and not limpia.startswith("//") else linea)
            continue

        # ----- declaraciones de nivel superior -----
        m = re.match(r"type\s+(%s)\s+struct\s*\{(.*)$" % IDENT, limpia)
        if m:
            indice_struct = len(salida)
            if m.group(2).strip().startswith("}"):
                salida.append("struct %s {};" % m.group(1))
            else:
                salida.append("struct %s {" % m.group(1))
                estado = "struct"
            continue
        m = re.match(r"type\s+(%s)\s+interface\s*\{(.*)$" % IDENT, limpia)
        if m:
            indice_struct = len(salida)
            if m.group(2).strip().startswith("}"):
                salida.append("class %s {};" % m.group(1))
            else:
                salida.append("class %s { public:" % m.group(1))
                estado = "interface"
            continue
        m = re.match(r"type\s+(%s)\s*=?\s*(.+)$" % IDENT, limpia)
        if m:
            salida.append("typedef %s %s;" % (convertir_tipo(m.group(2).split("//")[0]), m.group(1)))
            continue
        m = re.match(r"(const|var)\s*\($", limpia)
        if m:
            es_constante = m.group(1) == "const"
            estado = "bloque_var"
            salida.append("")
            continue
        m = re.match(r"(const|var)\s+(.+)$", limpia)
        if m:
            declaracion = declaracion_variable(m.group(2), m.group(1) == "const")
            salida.append(declaracion)
            sin = quitar_cadenas_y_comentarios(limpia)
            profundidad = sin.count("{") - sin.count("}")
            if profundidad > 0:
                estado = "cuerpo"
            continue
        if limpia.startswith("func"):
            sin = quitar_cadenas_y_comentarios(limpia)
            if sin.count("(") == sin.count(")") and sin.rstrip().endswith("{"):
                salida.append(emitir_funcion(limpia))
                profundidad = sin.count("{") - sin.count("}")
                estado = "cuerpo" if profundidad > 0 else None
            else:
                indice_struct = len(salida)
                salida.append("")
                cabecera = [limpia]
                estado = "cabecera"
            continue

        # Cualquier otra línea (comentarios //, líneas vacías) se conserva.
        salida.append(linea if limpia.startswith("//") or not limpia else "")

    if paquete_abierto:
        salida[-1] = salida[-1] + " } }"
    sys.stdout.write("\n".join(salida))


if __name__ == "__main__":
    filtrar(sys.argv[1])
