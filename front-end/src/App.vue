<script setup lang="ts">
import ArcadeButton from './components/ArcadeButton.vue'
</script>

<template>
  <div class="flex items-center justify-center min-h-screen gap-8">
    <ArcadeButton>
      <template #icon>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="size-6">
          <path d="M12 2L9.5 9.5H2l6 4.5L5.5 22 12 17l6.5 5-2.5-8 6-4.5h-7.5z" />
        </svg>
      </template>
    </ArcadeButton>

    <ArcadeButton label="Play" />

    <ArcadeButton label="Settings">
      <template #icon>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="size-6">
          <path d="M12 15.5A3.5 3.5 0 1 0 12 8.5a3.5 3.5 0 0 0 0 7z" />
          <path d="M19.43 12.98c.04-.32.07-.64.07-.98s-.03-.66-.07-.98l2.11-1.65a.5.5 0 0 0 .12-.64l-2-3.46a.5.5 0 0 0-.61-.22l-2.49 1a7.04 7.04 0 0 0-1.69-.98l-.38-2.65A.49.49 0 0 0 14 2h-4a.49.49 0 0 0-.49.42l-.38 2.65c-.61.25-1.17.59-1.69.98l-2.49-1a.5.5 0 0 0-.61.22l-2 3.46a.49.49 0 0 0 .12.64l2.11 1.65c-.04.32-.07.65-.07.98s.03.66.07.98l-2.11 1.65a.5.5 0 0 0-.12.64l2 3.46c.12.22.39.3.61.22l2.49-1c.52.4 1.08.73 1.69.98l.38 2.65c.05.24.26.42.49.42h4c.24 0 .44-.18.49-.42l.38-2.65c.61-.25 1.17-.59 1.69-.98l2.49 1c.22.08.49 0 .61-.22l2-3.46a.5.5 0 0 0-.12-.64l-2.11-1.65z" />
        </svg>
      </template>
    </ArcadeButton>
  </div>
</template>
