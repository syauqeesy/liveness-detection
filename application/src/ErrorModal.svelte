<script lang="ts">
  export let errorMessage: string | null

  export let onClose: () => void;
</script>

<div
  class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-5"
>
  <div class="w-full max-w-md rounded-2xl bg-white p-6 shadow-xl">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-2xl font-bold">Error</h2>
    </div>

    <div class="flex flex-col gap-12">
        <p class="text-center">
          {errorMessage ?? ''}
        </p>

      <button
        type="button"
        on:click={onClose}
        class="w-full rounded-xl border-3 p-3 font-semibold cursor-pointer hover:bg-black hover:text-white"
      >
        Close
      </button>
    </div>
  </div>
</div>

