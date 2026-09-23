<script lang="ts">
  export let result: {
    live: number;
    spoof: number;
    result: "Live" | "Spoof";
  };

  export let onClose: () => void;
</script>

<div
  class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-5"
>
  <div class="w-full max-w-md rounded-2xl bg-white p-6 shadow-xl">
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-2xl font-bold">Liveness Result</h2>

      <button
        type="button"
        on:click={onClose}
        class="text-2xl font-bold cursor-pointer"
      >
        x
      </button>
    </div>

    <div class="flex flex-col gap-4">
      <div class="text-center">
        <p class="text-sm text-gray-500">Result</p>

        <p class="text-4xl font-bold">
          {result.result}
        </p>
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div class="rounded-xl border p-4 text-center">
          <p class="text-sm text-gray-500">Live</p>

          <p class="text-2xl font-semibold">
            {(result.live * 100).toFixed(2)}%
          </p>
        </div>

        <div class="rounded-xl border p-4 text-center">
          <p class="text-sm text-gray-500">Spoof</p>

          <p class="text-2xl font-semibold">
            {(result.spoof * 100).toFixed(2)}%
          </p>
        </div>
      </div>

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

