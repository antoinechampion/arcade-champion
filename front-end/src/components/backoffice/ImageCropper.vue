<script setup lang="ts">
import { ref, computed, watch, onUnmounted } from 'vue'

const props = withDefaults(defineProps<{
  frameWidth: number
  frameHeight: number
  url: string
  outputScale?: number
}>(), { outputScale: 1 })

// In dev, external images must be proxied through localhost so the canvas can
// read their pixels (cross-origin taints the canvas). In production (Tauri),
// web security is disabled so no proxy is needed. Local paths (/images/...) are
// already same-origin and don't need proxying.
const proxiedUrl = computed(() => {
  const url = props.url
  if (!import.meta.env.DEV || !url.startsWith('http')) return url
  return `/api/local-dev-image-proxy?url=${encodeURIComponent(url)}`
})

const emit = defineEmits<{
  cropped: [dataUrl: string]
}>()

const imgRef = ref<HTMLImageElement>()

const loaded = ref(false)
const error = ref(false)

const scale = ref(1)
const translateX = ref(0)
const translateY = ref(0)

let imgNaturalWidth = 0
let imgNaturalHeight = 0

let dragging = false
let dragStartX = 0
let dragStartY = 0
let dragStartTX = 0
let dragStartTY = 0

function onLoad() {
  const img = imgRef.value!
  imgNaturalWidth = img.naturalWidth
  imgNaturalHeight = img.naturalHeight
  loaded.value = true
  error.value = false

  const scaleX = props.frameWidth / imgNaturalWidth
  const scaleY = props.frameHeight / imgNaturalHeight
  scale.value = Math.max(scaleX, scaleY)
  translateX.value = 0
  translateY.value = 0

  emitCrop()
}

function onError() {
  error.value = true
  loaded.value = false
}

function onWheel(e: WheelEvent) {
  e.preventDefault()
  const delta = e.deltaY > 0 ? 0.95 : 1.05
  const minScale = Math.max(props.frameWidth / imgNaturalWidth, props.frameHeight / imgNaturalHeight)
  scale.value = Math.max(minScale, scale.value * delta)
  clampTranslation()
  emitCrop()
}

function onPointerDown(e: PointerEvent) {
  dragging = true
  dragStartX = e.clientX
  dragStartY = e.clientY
  dragStartTX = translateX.value
  dragStartTY = translateY.value
  ;(e.target as HTMLElement).setPointerCapture(e.pointerId)
}

function onPointerMove(e: PointerEvent) {
  if (!dragging) return
  translateX.value = dragStartTX + (e.clientX - dragStartX)
  translateY.value = dragStartTY + (e.clientY - dragStartY)
  clampTranslation()
}

function onPointerUp() {
  if (!dragging) return
  dragging = false
  emitCrop()
}

function clampTranslation() {
  const scaledW = imgNaturalWidth * scale.value
  const scaledH = imgNaturalHeight * scale.value
  const maxX = (scaledW - props.frameWidth) / 2
  const maxY = (scaledH - props.frameHeight) / 2
  translateX.value = Math.max(-maxX, Math.min(maxX, translateX.value))
  translateY.value = Math.max(-maxY, Math.min(maxY, translateY.value))
}

function emitCrop() {
  const s = props.outputScale
  const canvas = document.createElement('canvas')
  canvas.width = props.frameWidth * s
  canvas.height = props.frameHeight * s
  const ctx = canvas.getContext('2d')!

  const scaledW = imgNaturalWidth * scale.value * s
  const scaledH = imgNaturalHeight * scale.value * s

  const drawX = (canvas.width - scaledW) / 2 + translateX.value * s
  const drawY = (canvas.height - scaledH) / 2 + translateY.value * s

  ctx.drawImage(imgRef.value!, drawX, drawY, scaledW, scaledH)
  emit('cropped', canvas.toDataURL('image/jpeg', 0.85))
}

watch(() => props.url, () => {
  loaded.value = false
  error.value = false
  scale.value = 1
  translateX.value = 0
  translateY.value = 0
})

onUnmounted(() => { dragging = false })
</script>

<template>
  <div class="cropper-wrapper">
    <div
      class="cropper-frame"
      :style="{ width: frameWidth + 'px', height: frameHeight + 'px' }"
      @wheel="onWheel"
      @pointerdown="onPointerDown"
      @pointermove="onPointerMove"
      @pointerup="onPointerUp"
      @pointercancel="onPointerUp"
    >
      <img
        v-show="loaded"
        ref="imgRef"
        :src="proxiedUrl"
        class="cropper-img"
        :style="{
          transform: `translate(-50%, -50%) translate(${translateX}px, ${translateY}px) scale(${scale})`,
        }"
        @load="onLoad"
        @error="onError"
      />
      <span v-if="!loaded && !error" class="placeholder">Loading…</span>
      <span v-if="error" class="placeholder error">Failed to load image</span>
    </div>
    <span class="hint">Scroll to zoom, drag to pan</span>
  </div>
</template>

<style scoped>
.cropper-wrapper {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}

.cropper-frame {
  position: relative;
  overflow: hidden;
  border-radius: 6px;
  border: 1px solid rgba(255, 255, 255, 0.15);
  background: rgba(0, 0, 0, 0.4);
  cursor: grab;
  touch-action: none;
  user-select: none;
}

.cropper-frame:active {
  cursor: grabbing;
}

.cropper-img {
  position: absolute;
  top: 50%;
  left: 50%;
  transform-origin: center;
  pointer-events: none;
  max-width: none;
}

.placeholder {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  opacity: 0.4;
}

.placeholder.error {
  color: #f87171;
  opacity: 0.8;
}

.hint {
  font-size: 0.6875rem;
  opacity: 0.4;
}
</style>
