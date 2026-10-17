ENTRY(_start)

SECTIONS
{
    .text :
    {
        . = ALIGN(16);
        KEEP(*(.init))
        . = ALIGN(16);
        *(.text)
        *(.text.*)
        *(.rodata)
        *(.rodata.*)
        . = ALIGN(16);
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

    /DISCARD/ :
    {
        *(.eh_frame)
    }
}

/* For the memory allocator. */
_heap_start = ALIGN(_ebss, 16);
_heap_end = ORIGIN(RAM) + LENGTH(RAM);
_globals_start = _sdata;
_globals_end = _ebss;
