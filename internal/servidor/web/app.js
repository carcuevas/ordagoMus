// Órdago · Mus en el navegador: un xterm.js conectado por websocket a una
// partida que corre en el servidor. Además emula lo justo del protocolo de
// gráficos de kitty (imágenes con marcadores Unicode) para ver la baraja
// ilustrada y las caras, y toca los sonidos que manda la partida.
"use strict";

const FONDO = "#1a1b26";

const term = new Terminal({
  fontFamily: '"JetBrainsMono Nerd Font", "DejaVu Sans Mono", Menlo, Consolas, monospace',
  fontSize: 15,
  cursorBlink: false,
  theme: { background: FONDO },
  allowProposedApi: true,
});
const fit = new FitAddon.FitAddon();
term.loadAddon(fit);
term.open(document.getElementById("terminal"));

// La mesa necesita unas 50 filas: se achica la letra hasta que quepan (sin
// bajar de 11 px), y se vuelve a crecer si hay sitio.
const FILAS_MESA = 50;
function ajustarLetra() {
  for (let px = 16; px >= 11; px--) {
    term.options.fontSize = px;
    fit.fit();
    if (term.rows >= FILAS_MESA) return;
  }
}
ajustarLetra();
term.focus();

// ───────────── gráficos de kitty ─────────────

// Diacríticos de kitty (rowcolumn-diacritics.txt): fila, columna y, solo en
// el modo web, la colocación.
const DIACRITICOS = [
  0x0305, 0x030d, 0x030e, 0x0310, 0x0312, 0x033d, 0x033e, 0x033f,
  0x0346, 0x034a, 0x034b, 0x034c, 0x0350, 0x0351, 0x0352, 0x0357,
  0x035b, 0x0363, 0x0364, 0x0365, 0x0366, 0x0367, 0x0368, 0x0369,
  0x036a, 0x036b, 0x036c, 0x036d, 0x036e, 0x036f,
];
const indiceDiacritico = new Map(DIACRITICOS.map((d, i) => [d, i]));

const imagenes = new Map();    // id → HTMLImageElement
const colocaciones = new Map(); // "id:pid" → {c, r}
let transmision = null;        // {id, partes[]} mientras llega una imagen por trozos

function comandoKitty(texto) {
  // texto es lo que va entre ESC _ y ESC \, empezando por "G".
  if (texto[0] !== "G") return;
  const pc = texto.indexOf(";");
  const claves = {};
  for (const kv of (pc < 0 ? texto.slice(1) : texto.slice(1, pc)).split(",")) {
    const i = kv.indexOf("=");
    if (i > 0) claves[kv.slice(0, i)] = kv.slice(i + 1);
  }
  const datos = pc < 0 ? "" : texto.slice(pc + 1);
  const a = claves.a || (transmision ? "t" : "");
  if (a === "t") {
    if (claves.i !== undefined) transmision = { id: Number(claves.i), partes: [] };
    if (!transmision) return;
    transmision.partes.push(datos);
    if (claves.m !== "1") {
      const { id, partes } = transmision;
      transmision = null;
      const img = new Image();
      img.onload = () => pintarTodo();
      img.src = "data:image/png;base64," + partes.join("");
      imagenes.set(id, img);
    }
  } else if (a === "p") {
    colocaciones.set(claves.i + ":" + claves.p, { c: Number(claves.c), r: Number(claves.r) });
  } else if (a === "d" && claves.i !== undefined) {
    imagenes.delete(Number(claves.i));
  }
}

// Filtro de la salida: quita los comandos APC (ESC _ … ESC \) antes de pasarla
// a xterm.js. Trabaja con bytes porque un comando puede venir partido.
// También cambia el marcador U+10EEEE (F4 8E BB AE en UTF-8) por un espacio:
// ninguna fuente lo tiene y sus cajitas asomarían por los lados de la imagen.
// Los diacríticos se quedan pegados al espacio, que es lo que se lee al pintar.
const ESC = 0x1b, GUION_BAJO = 0x5f, BARRA = 0x5c;
const MARCADOR_UTF8 = [0xf4, 0x8e, 0xbb, 0xae];
let estado = 0; // 0 normal · 1 tras ESC · 2 dentro de APC · 3 ESC dentro de APC · 4 posible marcador
let apc = [];
let marcador = 0; // bytes del marcador ya vistos
const decodificador = new TextDecoder();

function filtrar(bytes) {
  const salida = new Uint8Array(bytes.length + 8);
  let n = 0;
  for (const b of bytes) {
    if (estado === 4) {
      if (b === MARCADOR_UTF8[marcador]) {
        if (++marcador === 4) { salida[n++] = 0x20; estado = 0; }
        continue;
      }
      for (let i = 0; i < marcador; i++) salida[n++] = MARCADOR_UTF8[i];
      estado = 0;
    }
    switch (estado) {
      case 0:
        if (b === ESC) estado = 1;
        else if (b === MARCADOR_UTF8[0]) { estado = 4; marcador = 1; }
        else salida[n++] = b;
        break;
      case 1:
        if (b === GUION_BAJO) { estado = 2; apc = []; }
        else { salida[n++] = ESC; if (b === ESC) break; salida[n++] = b; estado = 0; }
        break;
      case 2:
        if (b === ESC) estado = 3; else apc.push(b);
        break;
      case 3:
        if (b === BARRA) {
          estado = 0;
          comandoKitty(decodificador.decode(new Uint8Array(apc)));
          apc = [];
        } else { apc.push(ESC, b); estado = 2; }
        break;
    }
  }
  return salida.subarray(0, n);
}

// Capa de imágenes encima del texto: allí donde hay un marcador se pinta su
// trozo de imagen.
const pantalla = term.element.querySelector(".xterm-screen");
const lienzo = document.createElement("canvas");
lienzo.id = "imagenes";
pantalla.appendChild(lienzo);
const ctx = lienzo.getContext("2d");

function celda() {
  const d = term._core && term._core._renderService && term._core._renderService.dimensions;
  if (d && d.css && d.css.cell && d.css.cell.width) return { w: d.css.cell.width, h: d.css.cell.height };
  return { w: pantalla.clientWidth / term.cols, h: pantalla.clientHeight / term.rows };
}

function ajustarLienzo() {
  const { w, h } = celda();
  const dpr = window.devicePixelRatio || 1;
  const ancho = Math.round(w * term.cols), alto = Math.round(h * term.rows);
  if (lienzo.width !== Math.round(ancho * dpr) || lienzo.height !== Math.round(alto * dpr)) {
    lienzo.width = Math.round(ancho * dpr);
    lienzo.height = Math.round(alto * dpr);
    lienzo.style.width = ancho + "px";
    lienzo.style.height = alto + "px";
  }
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
  ctx.imageSmoothingQuality = "high";
}

// Se recorre la pantalla buscando marcadores; de cada uno se deduce dónde
// empieza su imagen (restando su fila y columna) y cada imagen se pinta entera
// una sola vez: pintarla celda a celda deja rayas entre trozos.
let celdaReusable;
function pintarTodo() {
  ajustarLienzo();
  const { w, h } = celda();
  const buf = term.buffer.active;
  ctx.clearRect(0, 0, term.cols * w + 1, term.rows * h + 1);
  const hechas = new Set();
  for (let y = 0; y < term.rows; y++) {
    const linea = buf.getLine(buf.viewportY + y);
    if (!linea) continue;
    for (let x = 0; x < term.cols; x++) {
      celdaReusable = linea.getCell(x, celdaReusable);
      if (!celdaReusable) continue;
      const ch = celdaReusable.getChars();
      if (!ch || ch.length < 3 || ch[0] !== " " || !celdaReusable.isFgRGB()) continue;
      const cps = Array.from(ch).slice(1).map((c) => c.codePointAt(0));
      const fila = indiceDiacritico.get(cps[0]);
      const col = indiceDiacritico.get(cps[1]);
      const pid = indiceDiacritico.get(cps[2]);
      if (fila === undefined || col === undefined || pid === undefined) continue;
      const id = celdaReusable.getFgColor();
      const x0 = x - col, y0 = y - fila;
      const clave = id + ":" + pid + "@" + x0 + "," + y0;
      if (hechas.has(clave)) continue;
      hechas.add(clave);
      const img = imagenes.get(id);
      const sitio = colocaciones.get(id + ":" + pid);
      if (!img || !img.complete || !img.naturalWidth || !sitio) continue;
      const px = Math.round(x0 * w), py = Math.round(y0 * h);
      const ancho = Math.round((x0 + sitio.c) * w) - px, alto = Math.round((y0 + sitio.r) * h) - py;
      ctx.fillStyle = FONDO;
      ctx.fillRect(px, py, ancho, alto);
      ctx.drawImage(img, px, py, ancho, alto);
    }
  }
}

let pendiente = false;
function pintarLuego() {
  if (pendiente) return;
  pendiente = true;
  requestAnimationFrame(() => { pendiente = false; pintarTodo(); });
}
term.onRender(() => pintarLuego());
term.onResize(() => pintarLuego());

// ───────────── sonidos ─────────────

const sonidos = {};
for (const n of ["reparto", "corte", "lance", "envite", "ordago", "recuento"]) {
  const a = new Audio("sonidos/" + n + ".wav");
  a.preload = "auto";
  sonidos[n] = a;
}
function tocar(n) {
  const a = sonidos[n];
  if (!a) return;
  const b = a.cloneNode();
  b.play().catch(() => {}); // el navegador no deja sonar antes de pulsar una tecla
}

// ───────────── conexión ─────────────

const TEXTOS = {
  es: ["La partida ha terminado.", "Volver a jugar"],
  en: ["The game is over.", "Play again"],
  cs: ["Hra skončila.", "Hrát znovu"],
  eu: ["Partida amaitu da.", "Berriro jokatu"],
};
const texto = TEXTOS[(navigator.language || "es").slice(0, 2)] || TEXTOS.es;
document.getElementById("reconectar").textContent = texto[1];

const CLAVE_AJUSTES = "ordago.ajustes";
let ws = null;

function avisar(texto) {
  document.getElementById("aviso-texto").textContent = texto;
  document.getElementById("aviso").hidden = false;
}

function tam(tipo) {
  const { w, h } = celda();
  return { t: tipo, cols: term.cols, rows: term.rows, cw: w, ch: h };
}

function conectar() {
  document.getElementById("aviso").hidden = true;
  imagenes.clear();
  colocaciones.clear();
  estado = 0;
  term.reset();
  const url = (location.protocol === "https:" ? "wss://" : "ws://") + location.host + "/ws";
  ws = new WebSocket(url);
  ws.binaryType = "arraybuffer";
  ws.onopen = () => {
    const ini = tam("inicio");
    ini.ajustes = localStorage.getItem(CLAVE_AJUSTES) || "";
    ws.send(JSON.stringify(ini));
  };
  ws.onmessage = (ev) => {
    if (typeof ev.data === "string") {
      let m;
      try { m = JSON.parse(ev.data); } catch { return; }
      if (m.t === "sonido") tocar(m.n);
      else if (m.t === "ajustes" && typeof m.d === "string") localStorage.setItem(CLAVE_AJUSTES, m.d);
      return;
    }
    const datos = filtrar(new Uint8Array(ev.data));
    if (datos.length) term.write(datos);
  };
  ws.onclose = (ev) => {
    avisar(ev.reason || texto[0]);
  };
}

const codificador = new TextEncoder();
term.onData((d) => { if (ws && ws.readyState === WebSocket.OPEN) ws.send(codificador.encode(d)); });
term.onBinary((d) => {
  if (!ws || ws.readyState !== WebSocket.OPEN) return;
  const b = new Uint8Array(d.length);
  for (let i = 0; i < d.length; i++) b[i] = d.charCodeAt(i) & 0xff;
  ws.send(b);
});

let espera;
window.addEventListener("resize", () => {
  clearTimeout(espera);
  espera = setTimeout(() => {
    ajustarLetra();
    if (ws && ws.readyState === WebSocket.OPEN) ws.send(JSON.stringify(tam("tam")));
    pintarTodo();
  }, 100);
});

document.getElementById("reconectar").addEventListener("click", () => { conectar(); term.focus(); });
conectar();
