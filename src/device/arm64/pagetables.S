// Page table storage for the MMU. One L0 table, one L1 table,
// two L2 tables, and 32 L3 tables of 512 entries each, 4096 bytes
// per table. The linker script places this section 4096 byte
// aligned, outside .bss, so the tables survive BSS zeroing. The
// L3 tables cover 64M of RAM at 4K pages. See ARM ARM, table D8-1
// for the descriptor format.
.section .pagetables, "aw", @nobits
.balign 4096
.global _pagetables_storage
_pagetables_storage:
    .zero 36 * 4096
