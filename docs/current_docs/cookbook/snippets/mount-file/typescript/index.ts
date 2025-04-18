import { dag, Container, File, object, func } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Return a container with a mounted file
   */
  @func()
  async copyFile(
    /**
     * Source file
     */
    f: File,
  ): Promise<Container> {
    const name = await f.name()
    return dag
      .container()
      .from("alpine:latest")
      .withMountedFile(`/src/${name}`, f)
  }
}
