PORTB = $6000
PORTA = $6001
DDRB = $6002
DDRA = $6003

message = $0001
stringPointer = $007f

value = $0200
mod10 = $0202

E  = %10000000
RW = %01000000
RS = %00100000


  .org $8000

reset:
  lda #%11111111 ; Set all pins on port B to output
  sta DDRB
  lda #%11100000 ; Set top 3 pins on port A to output
  sta DDRA

  lda #%00111000 ; Set 8-bit mode; 2-line display; 5x8 font
  jsr lcd_instruction
  lda #%00001110 ; Display on; cursor on; blink off
  jsr lcd_instruction
  lda #%00000110 ; Increment and shift cursor; don't shift display
  jsr lcd_instruction
  lda #$00000001 ; Clear display
  jsr lcd_instruction

  jmp loop

string: .asciiz "Running Program .."
printString:
  lda #$00000001 ; Clear display
  jsr lcd_instruction
  ldx #0
printS:
  lda string,x
  beq return_str
  jsr print_char
  inx
  jmp printS
return_str:
  rts


printNumber:
  lda #$00000001 ; Clear display
  jsr lcd_instruction
  lda #0
  sta message
divide :
  lda #0
  sta mod10
  sta mod10 + 1
  clc 
  ldx #16
divloop:
  rol value
  rol value + 1
  rol mod10
  rol mod10 + 1
  sec 
  lda mod10
  sbc #10
  tay 
  lda mod10 + 1
  sbc #0
  bcc ingnore_result
  sty mod10
  sta mod10 + 1
ingnore_result: 
  dex 
  bne divloop
  rol value
  rol value + 1
  lda mod10
  clc
  adc #"0"
  jsr push_char
  lda value
  ora value + 1
  bne divide
  ldx #0
print:
  lda message,x
  beq return_bin
  jsr print_char
  inx
  jmp print
return_bin:
  rts

loop:
  jmp loop

push_char:
  pha
  ldy #0
char_loop:
  lda message,y
  tax 
  pla 
  sta message,y
  iny 
  txa 
  pha
  bne char_loop
  pla
  sta message,y
  rts

lcd_wait:
  pha
  lda #%00000000  ; Port B is input
  sta DDRB
lcdbusy:
  lda #RW
  sta PORTA
  lda #(RW | E)
  sta PORTA
  lda PORTB
  and #%10000000
  bne lcdbusy
  lda #RW
  sta PORTA
  lda #%11111111  ; Port B is output
  sta DDRB
  pla
  rts

lcd_instruction:
  jsr lcd_wait
  sta PORTB
  lda #0         ; Clear RS/RW/E bits
  sta PORTA
  lda #E         ; Set E bit to send instruction
  sta PORTA
  lda #0         ; Clear RS/RW/E bits
  sta PORTA
  rts

print_char:
  jsr lcd_wait
  sta PORTB
  lda #RS         ; Set RS; Clear RW/E bits
  sta PORTA
  lda #(RS | E)   ; Set E bit to send instruction
  sta PORTA
  lda #RS         ; Clear E bits
  sta PORTA
  rts

delay:
  ldy #$ff
  ldx #$ff
delayO:
  dex
  bne delayO
  dey 
  bne delayO
  ldy #$ff
  ldx #$ff
delayT:
  dex
  bne delayT
  dey 
  bne delayT
  rts
  

  .org $fffc
  .word reset
  .word $0000