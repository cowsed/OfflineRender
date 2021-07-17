# Offline Renderer


A path tracer based off of (Raytracing in One Weekend)[https://raytracing.github.io/books/RayTracingInOneWeekend.html#diffusematerials/truelambertianreflection)

Problems encountered and solved:
- Parralellizing the render
  - Fairly easy fix in go by splitting the image into chunks and assigning chunks to different Goroutines
- Fixing a locking random number generator when parallelized
  - Presumably to insure no duplicate numbers can be duplicated by the default source, the pseudorandom number generator can only provide one number to one goroutine at one time
  - This was fixed by giving each goroutine its own random source and passing references to it through all functions that needed it
