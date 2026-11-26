ENTRY(_start)

SECTIONS
{
    .vectors :
    {
        . = ALIGN(2048);
        KEEP(*(.vectors))
    } >FLASH_TEXT

    .text :
    {
        . = ALIGN(16);
        _stext = .;        /* used by MMU setup: start of read-only code */
        KEEP(*(.init))
        . = ALIGN(16);
        *(.text)
        *(.text.*)
        *(.rodata)
        *(.rodata.*)
        . = ALIGN(16);
        _etext = .;        /* used by MMU setup: end of read-only code */
        . = ALIGN(4096);   /* code/data page boundary for MMU permissions */
    } >FLASH_TEXT

    /* Put the stack at the bottom of RAM, so that the application will
     * crash on stack overflow instead of silently corrupting memory.
     * See: http://blog.japaric.io/stack-overflow-protection/ */
    .stack (NOLOAD) :
    {
        . = ALIGN(16);
        . += _stack_size;
        _stack_top = .;
    } >RAM

    /* Start address (in flash) of .data, used by startup code. */
    _sidata = LOADADDR(.data);

    /* Globals with initial value */
    .data :
    {
        . = ALIGN(8);
        _sdata = .;        /* used by startup code */
        *(.sdata)
        *(.data .data.*)
        *(.got)
        *(.got.plt)
        . = ALIGN(8);
        _edata = .;        /* used by startup code */
    } >RAM AT>FLASH_TEXT

    /* Zero-initialized globals  */
    .bss :
    {
        . = ALIGN(8);
        _sbss = .;         /* used by startup code */
        *(.sbss)
        *(.bss .bss.*)
        *(COMMON)
        . = ALIGN(8);
        _ebss = .;         /* used by startup code */
    } >RAM

    /* Page tables for the MMU. Outside .bss so startup BSS zeroing
     * cannot erase tables built before main, and 4096 byte aligned
     * as TTBR0_EL1 requires. Sized for one L0, one L1, two L2, and
     * 32 L3 tables, enough for 64M of RAM at 4K pages. The heap
     * starts after the tables. */
    .pagetables (NOLOAD) :
    {
        . = ALIGN(4096);
        _pagetables = .;   /* used by MMU setup */
        KEEP(*(.pagetables))
        . = ALIGN(4096);
        _pagetables_end = .;
    } >RAM

    /DISCARD/ :
    {
        *(.eh_frame)
    }
}

/* Start of RAM, used by MMU setup as the identity map base. */
_ram_start = ORIGIN(RAM);

/* For the memory allocator. Page tables sit after .bss. */
_heap_start = ALIGN(_pagetables_end, 16);
_heap_end = ORIGIN(RAM) + LENGTH(RAM);
_globals_start = _sdata;
_globals_end = _ebss;
