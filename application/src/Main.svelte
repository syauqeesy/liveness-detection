<script lang="ts">
  import { onMount } from "svelte";
  import type {
    InferenceServiceAdapter,
    InferenceResult,
  } from "./service/main";
  import ManagedService from "./service/managed-service";
  import SelfManagedService from "./service/self-managed-service";
  import OnDeviceService from "./service/on-device-service";
  import ResultModal from "./ResultModal.svelte";
  import Loading from "./Loading.svelte";

  let IsLoading = true;
  let VideoElement: HTMLVideoElement;
  let Stream: MediaStream | null = null;
  let StreamHeight = 0;
  let inferenceFinished = false;
  let inferenceResult: InferenceResult | null = null;

  const managedService = new ManagedService();
  const selfManagedService = new SelfManagedService();
  const onDeviceService = new OnDeviceService();

  async function startCamera() {
    IsLoading = true;

    try {
      Stream = await navigator.mediaDevices.getUserMedia({
        video: {
          width: { ideal: 1280 },
          height: { ideal: 720 },
        },
        audio: false,
      });

      VideoElement.srcObject = Stream;

      await new Promise<void>((resolve, reject) => {
        const onReady = async () => {
          try {
            await VideoElement.play();

            if (
              VideoElement.videoWidth === 0 ||
              VideoElement.videoHeight === 0
            ) {
              reject(new Error("Camera video dimensions are not available"));
              return;
            }

            StreamHeight = VideoElement.getBoundingClientRect().height;

            resolve();
          } catch (error) {
            reject(error);
          }
        };

        const onError = () => {
          reject(new Error("Failed to start camera video"));
        };

        VideoElement.addEventListener("loadedmetadata", onReady, {
          once: true,
        });

        VideoElement.addEventListener("error", onError, {
          once: true,
        });
      });

      IsLoading = false;
    } catch (error: unknown) {
      if (error instanceof Error) {
        console.error(error.message);
      } else {
        console.error("Failed to start camera");
      }

      IsLoading = false;
    }
  }

  async function takeCroppedPicture(): Promise<Blob> {
    if (!VideoElement) {
      throw new Error("Video element is not available");
    }

    if (VideoElement.videoWidth === 0 || VideoElement.videoHeight === 0) {
      throw new Error("Camera is not ready");
    }

    const videoRect = VideoElement.getBoundingClientRect();

    const cropDisplaySize = 448;
    const outputSize = 128;

    const cropWidth =
      cropDisplaySize * (VideoElement.videoWidth / videoRect.width);
    const cropHeight =
      cropDisplaySize * (VideoElement.videoHeight / videoRect.height);

    const sourceX = (VideoElement.videoWidth - cropWidth) / 2;
    const sourceY = (VideoElement.videoHeight - cropHeight) / 2;

    const canvas = document.createElement("canvas");

    canvas.width = outputSize;
    canvas.height = outputSize;

    const context = canvas.getContext("2d");

    if (!context) {
      throw new Error("Could not get canvas context");
    }

    context.drawImage(
      VideoElement,

      sourceX,
      sourceY,
      cropWidth,
      cropHeight,

      0,
      0,
      outputSize,
      outputSize,
    );

    const blob = await new Promise<Blob | null>((resolve) => {
      canvas.toBlob(resolve, "image/jpeg", 0.9);
    });

    if (!blob) {
      throw new Error("Failed to create image");
    }

    return blob;
  }

  async function execute(
    mode: "managed_service" | "self_managed_service" | "on_device_service",
  ) {
    IsLoading = true;

    try {
      const imageBlob = await takeCroppedPicture();

      let InferenceService: InferenceServiceAdapter;

      switch (mode) {
        case "managed_service":
          InferenceService = managedService;
          break;
        case "self_managed_service":
          InferenceService = selfManagedService;
          break;
        case "on_device_service":
          InferenceService = onDeviceService;
          break;
      }

      inferenceResult = await InferenceService.Predict(imageBlob);
      inferenceFinished = true;
    } catch (error) {
      console.error("Liveness check failed:", error);
    } finally {
      IsLoading = false;
    }
  }

  function onCloseModal() {
    inferenceFinished = false;
    inferenceResult = null;
  }

  onMount(async () => {
    IsLoading = true;

    try {
      await onDeviceService.Initialize();
      await startCamera();
    } catch (error) {
      console.error("Initialization failed:", error);
      IsLoading = false;
    }
  });
</script>

{#if IsLoading}
  <Loading></Loading>
{/if}

<main class="flex flex-col gap-10 xl:p-10 md:p-5 p-2">
  <section id="camera-preview-container" class="flex flex-col items-center">
    <div
      class="xl:w-[1280px] w-full xl:h-[720px] min-h-12 rounded-xl border-solid border-black outline-3"
    >
      <video
        class={"w-full rounded-xl" +
          " " +
          (Stream === null ? "hidden" : "block")}
        bind:this={VideoElement}
        autoplay
        playsinline
        muted
      ></video>
    </div>

    {#if !IsLoading}
      <div
        style:top={`${StreamHeight > 0 ? (StreamHeight - 448) / 2 : 0}px`}
        class="absolute left-1/2 -translate-x-1/2 xl:w-[448px] xl:h-[448px] w-[224px] h-[224px] border-solid border-black outline-3"
      ></div>
    {/if}
  </section>

  <section id="camera-control-container" class="flex justify-center gap-5">
    <button
      on:click={() => execute("managed_service")}
      class="p-3 rounded-xl border-3 font-semibold cursor-pointer hover:bg-black hover:text-white xl:size-fit md:size-fit w-full"
      >Managed Service Liveness Check</button
    >
    <button
      on:click={() => execute("self_managed_service")}
      class="p-3 rounded-xl border-3 font-semibold cursor-pointer hover:bg-black hover:text-white xl:size-fit md:size-fit w-full"
      >Self Managed Service Liveness Check</button
    >
    <button
      on:click={() => execute("on_device_service")}
      class="p-3 rounded-xl border-3 font-semibold cursor-pointer hover:bg-black hover:text-white xl:size-fit md:size-fit w-full"
      >On Device Service Liveness Check</button
    >
  </section>
</main>

{#if inferenceFinished && inferenceResult}
  <ResultModal result={inferenceResult} onClose={onCloseModal}></ResultModal>
{/if}
