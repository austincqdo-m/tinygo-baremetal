/* Memory map for the QEMU virt machine.
 * RAM starts at 0x40000000. The size is kept below the 128M given to QEMU
 * so the heap stays small enough for fast tests. */
MEMORY
{
    RAM (rwx) : ORIGIN = 0x40000000, LENGTH = 64M
}

REGION_ALIAS("FLASH_TEXT", RAM)

_stack_size = 64K;

INCLUDE "targets/aarch64.ld"
