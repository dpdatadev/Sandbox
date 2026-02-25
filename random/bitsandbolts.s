	.arch armv8-a
	.file	"bitsandbolts.cpp"
	.text
#APP
	.globl _ZSt21ios_base_library_initv
	.section	.rodata
	.align	3
.LC0:
	.string	"success.."
	.align	3
.LC1:
	.string	"FAIL"
#NO_APP
	.text
	.align	2
	.global	_Z5test1v
	.type	_Z5test1v, %function
_Z5test1v:
.LFB2009:
	.cfi_startproc
	stp	x29, x30, [sp, -32]!
	.cfi_def_cfa_offset 32
	.cfi_offset 29, -32
	.cfi_offset 30, -24
	mov	x29, sp
	str	wzr, [sp, 28]
	ldr	w1, [sp, 28]
	adrp	x0, :got:_ZSt4cout;ldr	x0, [x0, :got_lo12:_ZSt4cout]
	bl	_ZNSolsEj
	adrp	x1, :got:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_;ldr	x1, [x1, :got_lo12:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_]
	bl	_ZNSolsEPFRSoS_E
	ldr	w0, [sp, 28]
	orr	w0, w0, 2
	str	w0, [sp, 28]
	ldr	w1, [sp, 28]
	adrp	x0, :got:_ZSt4cout;ldr	x0, [x0, :got_lo12:_ZSt4cout]
	bl	_ZNSolsEj
	adrp	x1, :got:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_;ldr	x1, [x1, :got_lo12:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_]
	bl	_ZNSolsEPFRSoS_E
	ldr	w0, [sp, 28]
	and	w0, w0, 2
	cmp	w0, 0
	beq	.L2
	adrp	x0, .LC0
	add	x1, x0, :lo12:.LC0
	adrp	x0, :got:_ZSt4cout;ldr	x0, [x0, :got_lo12:_ZSt4cout]
	bl	_ZStlsISt11char_traitsIcEERSt13basic_ostreamIcT_ES5_PKc
	adrp	x1, :got:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_;ldr	x1, [x1, :got_lo12:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_]
	bl	_ZNSolsEPFRSoS_E
	b	.L3
.L2:
	adrp	x0, .LC1
	add	x1, x0, :lo12:.LC1
	adrp	x0, :got:_ZSt4cerr;ldr	x0, [x0, :got_lo12:_ZSt4cerr]
	bl	_ZStlsISt11char_traitsIcEERSt13basic_ostreamIcT_ES5_PKc
	adrp	x1, :got:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_;ldr	x1, [x1, :got_lo12:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_]
	bl	_ZNSolsEPFRSoS_E
.L3:
	ldr	w0, [sp, 28]
	str	w0, [sp, 24]
	ldr	w0, [sp, 24]
	orr	w0, w0, 1
	str	w0, [sp, 24]
	ldr	w0, [sp, 24]
	and	w0, w0, 1
	cmp	w0, 0
	beq	.L4
	adrp	x0, .LC0
	add	x1, x0, :lo12:.LC0
	adrp	x0, :got:_ZSt4cout;ldr	x0, [x0, :got_lo12:_ZSt4cout]
	bl	_ZStlsISt11char_traitsIcEERSt13basic_ostreamIcT_ES5_PKc
	adrp	x1, :got:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_;ldr	x1, [x1, :got_lo12:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_]
	bl	_ZNSolsEPFRSoS_E
	b	.L6
.L4:
	adrp	x0, .LC1
	add	x1, x0, :lo12:.LC1
	adrp	x0, :got:_ZSt4cerr;ldr	x0, [x0, :got_lo12:_ZSt4cerr]
	bl	_ZStlsISt11char_traitsIcEERSt13basic_ostreamIcT_ES5_PKc
	adrp	x1, :got:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_;ldr	x1, [x1, :got_lo12:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_]
	bl	_ZNSolsEPFRSoS_E
.L6:
	nop
	ldp	x29, x30, [sp], 32
	.cfi_restore 30
	.cfi_restore 29
	.cfi_def_cfa_offset 0
	ret
	.cfi_endproc
.LFE2009:
	.size	_Z5test1v, .-_Z5test1v
	.align	2
	.global	_Z5test2v
	.type	_Z5test2v, %function
_Z5test2v:
.LFB2010:
	.cfi_startproc
	sub	sp, sp, #16
	.cfi_def_cfa_offset 16
	mov	w0, 5
	strh	w0, [sp, 14]
	ldrh	w0, [sp, 14]
	ubfiz	w0, w0, 2, 14
	strh	w0, [sp, 14]
	ldrsh	w0, [sp, 14]
	add	sp, sp, 16
	.cfi_def_cfa_offset 0
	ret
	.cfi_endproc
.LFE2010:
	.size	_Z5test2v, .-_Z5test2v
	.align	2
	.global	_Z5test3v
	.type	_Z5test3v, %function
_Z5test3v:
.LFB2011:
	.cfi_startproc
	sub	sp, sp, #16
	.cfi_def_cfa_offset 16
	mov	w0, 5
	str	w0, [sp, 12]
	ldr	w0, [sp, 12]
	lsl	w0, w0, 2
	str	w0, [sp, 12]
	ldr	w0, [sp, 12]
	add	sp, sp, 16
	.cfi_def_cfa_offset 0
	ret
	.cfi_endproc
.LFE2011:
	.size	_Z5test3v, .-_Z5test3v
	.align	2
	.global	_Z5test4v
	.type	_Z5test4v, %function
_Z5test4v:
.LFB2012:
	.cfi_startproc
	sub	sp, sp, #16
	.cfi_def_cfa_offset 16
	mov	w0, 48
	str	w0, [sp, 12]
	ldr	w0, [sp, 12]
	asr	w0, w0, 2
	str	w0, [sp, 12]
	ldr	w0, [sp, 12]
	add	sp, sp, 16
	.cfi_def_cfa_offset 0
	ret
	.cfi_endproc
.LFE2012:
	.size	_Z5test4v, .-_Z5test4v
	.section	.rodata
	.align	3
.LC2:
	.string	"Bitmasking Etc.,\n"
	.align	3
.LC3:
	.string	"The bit is set.."
	.align	3
.LC4:
	.string	"The bit is not set.."
	.text
	.align	2
	.global	main
	.type	main, %function
main:
.LFB2013:
	.cfi_startproc
	stp	x29, x30, [sp, -48]!
	.cfi_def_cfa_offset 48
	.cfi_offset 29, -48
	.cfi_offset 30, -40
	mov	x29, sp
	adrp	x0, .LC2
	add	x1, x0, :lo12:.LC2
	adrp	x0, :got:_ZSt4cout;ldr	x0, [x0, :got_lo12:_ZSt4cout]
	bl	_ZStlsISt11char_traitsIcEERSt13basic_ostreamIcT_ES5_PKc
	adrp	x1, :got:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_;ldr	x1, [x1, :got_lo12:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_]
	bl	_ZNSolsEPFRSoS_E
	str	wzr, [sp, 44]
	ldr	w1, [sp, 44]
	adrp	x0, :got:_ZSt4cout;ldr	x0, [x0, :got_lo12:_ZSt4cout]
	bl	_ZNSolsEj
	adrp	x1, :got:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_;ldr	x1, [x1, :got_lo12:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_]
	bl	_ZNSolsEPFRSoS_E
	ldr	w0, [sp, 44]
	orr	w0, w0, 1
	str	w0, [sp, 44]
	ldr	w1, [sp, 44]
	adrp	x0, :got:_ZSt4cout;ldr	x0, [x0, :got_lo12:_ZSt4cout]
	bl	_ZNSolsEj
	adrp	x1, :got:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_;ldr	x1, [x1, :got_lo12:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_]
	bl	_ZNSolsEPFRSoS_E
	str	wzr, [sp, 44]
	ldr	w1, [sp, 44]
	adrp	x0, :got:_ZSt4cout;ldr	x0, [x0, :got_lo12:_ZSt4cout]
	bl	_ZNSolsEj
	adrp	x1, :got:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_;ldr	x1, [x1, :got_lo12:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_]
	bl	_ZNSolsEPFRSoS_E
	mov	w0, 1
	str	w0, [sp, 40]
	ldr	w0, [sp, 40]
	lsl	w0, w0, 3
	str	w0, [sp, 40]
	ldr	w0, [sp, 44]
	orr	w0, w0, 8
	str	w0, [sp, 44]
	ldr	w1, [sp, 44]
	adrp	x0, :got:_ZSt4cout;ldr	x0, [x0, :got_lo12:_ZSt4cout]
	bl	_ZNSolsEj
	adrp	x1, :got:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_;ldr	x1, [x1, :got_lo12:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_]
	bl	_ZNSolsEPFRSoS_E
	ldr	w0, [sp, 44]
	eor	w0, w0, 8
	str	w0, [sp, 44]
	ldr	w1, [sp, 44]
	adrp	x0, :got:_ZSt4cout;ldr	x0, [x0, :got_lo12:_ZSt4cout]
	bl	_ZNSolsEj
	adrp	x1, :got:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_;ldr	x1, [x1, :got_lo12:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_]
	bl	_ZNSolsEPFRSoS_E
	ldr	w0, [sp, 44]
	eor	w0, w0, 8
	str	w0, [sp, 44]
	ldr	w0, [sp, 44]
	and	w0, w0, 8
	cmp	w0, 0
	beq	.L14
	adrp	x0, .LC3
	add	x1, x0, :lo12:.LC3
	adrp	x0, :got:_ZSt4cout;ldr	x0, [x0, :got_lo12:_ZSt4cout]
	bl	_ZStlsISt11char_traitsIcEERSt13basic_ostreamIcT_ES5_PKc
	adrp	x1, :got:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_;ldr	x1, [x1, :got_lo12:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_]
	bl	_ZNSolsEPFRSoS_E
	b	.L15
.L14:
	adrp	x0, .LC4
	add	x1, x0, :lo12:.LC4
	adrp	x0, :got:_ZSt4cout;ldr	x0, [x0, :got_lo12:_ZSt4cout]
	bl	_ZStlsISt11char_traitsIcEERSt13basic_ostreamIcT_ES5_PKc
	adrp	x1, :got:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_;ldr	x1, [x1, :got_lo12:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_]
	bl	_ZNSolsEPFRSoS_E
.L15:
	strb	wzr, [sp, 39]
	ldrb	w0, [sp, 39]
	orr	w0, w0, 1
	strb	w0, [sp, 39]
	ldrb	w1, [sp, 39]
	adrp	x0, :got:_ZSt4cout;ldr	x0, [x0, :got_lo12:_ZSt4cout]
	bl	_ZStlsISt11char_traitsIcEERSt13basic_ostreamIcT_ES5_h
	adrp	x1, :got:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_;ldr	x1, [x1, :got_lo12:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_]
	bl	_ZNSolsEPFRSoS_E
	ldrb	w0, [sp, 39]
	orr	w0, w0, 4
	strb	w0, [sp, 39]
	ldrb	w1, [sp, 39]
	adrp	x0, :got:_ZSt4cout;ldr	x0, [x0, :got_lo12:_ZSt4cout]
	bl	_ZStlsISt11char_traitsIcEERSt13basic_ostreamIcT_ES5_h
	adrp	x1, :got:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_;ldr	x1, [x1, :got_lo12:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_]
	bl	_ZNSolsEPFRSoS_E
	mov	w0, 214
	str	w0, [sp, 32]
	ldr	w0, [sp, 32]
	lsr	w0, w0, 5
	and	w0, w0, 7
	str	w0, [sp, 28]
	ldr	w1, [sp, 28]
	adrp	x0, :got:_ZSt4cout;ldr	x0, [x0, :got_lo12:_ZSt4cout]
	bl	_ZNSolsEj
	adrp	x1, :got:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_;ldr	x1, [x1, :got_lo12:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_]
	bl	_ZNSolsEPFRSoS_E
	bl	_Z5test1v
	bl	_Z5test2v
	sxth	w0, w0
	mov	w1, w0
	adrp	x0, :got:_ZSt4cout;ldr	x0, [x0, :got_lo12:_ZSt4cout]
	bl	_ZNSolsEs
	adrp	x1, :got:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_;ldr	x1, [x1, :got_lo12:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_]
	bl	_ZNSolsEPFRSoS_E
	bl	_Z5test3v
	mov	w1, w0
	adrp	x0, :got:_ZSt4cout;ldr	x0, [x0, :got_lo12:_ZSt4cout]
	bl	_ZNSolsEi
	adrp	x1, :got:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_;ldr	x1, [x1, :got_lo12:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_]
	bl	_ZNSolsEPFRSoS_E
	bl	_Z5test4v
	mov	w1, w0
	adrp	x0, :got:_ZSt4cout;ldr	x0, [x0, :got_lo12:_ZSt4cout]
	bl	_ZNSolsEi
	adrp	x1, :got:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_;ldr	x1, [x1, :got_lo12:_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_]
	bl	_ZNSolsEPFRSoS_E
	mov	w0, 0
	ldp	x29, x30, [sp], 48
	.cfi_restore 30
	.cfi_restore 29
	.cfi_def_cfa_offset 0
	ret
	.cfi_endproc
.LFE2013:
	.size	main, .-main
	.section	.rodata
	.type	_ZNSt8__detail30__integer_to_chars_is_unsignedIjEE, %object
	.size	_ZNSt8__detail30__integer_to_chars_is_unsignedIjEE, 1
_ZNSt8__detail30__integer_to_chars_is_unsignedIjEE:
	.byte	1
	.type	_ZNSt8__detail30__integer_to_chars_is_unsignedImEE, %object
	.size	_ZNSt8__detail30__integer_to_chars_is_unsignedImEE, 1
_ZNSt8__detail30__integer_to_chars_is_unsignedImEE:
	.byte	1
	.type	_ZNSt8__detail30__integer_to_chars_is_unsignedIyEE, %object
	.size	_ZNSt8__detail30__integer_to_chars_is_unsignedIyEE, 1
_ZNSt8__detail30__integer_to_chars_is_unsignedIyEE:
	.byte	1
	.ident	"GCC: (Debian 14.2.0-19) 14.2.0"
	.section	.note.GNU-stack,"",@progbits
