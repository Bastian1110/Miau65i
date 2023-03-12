# Miau6502 Improved 🧶

A programming language for Ben Eater's 6502, using the assembler syntaxis presented in Ben's videos

## About

The objective of this project, is to create a super-simple progrmamming language, that
its [Turing Complete](https://en.wikipedia.org/wiki/Turing_completeness). This laguange
will compile fot the 65C02 legacy processor, more precisely for the Ben Eater project.

#### **WARNING : Currently in development**

## Objectives

1. Typed language
2. Support functions for 65C22
3. Support functions for LCD Display
4. Implement an standard libary

## How it works ?

This compiler will implement the most basic and common algorithms that are
normally in a compiler. This to demonstrate the way a compiler work (and a little bit of a challenge for myself 😅). The main parts of the compiler are described below

## Lexer

For the implementation of the lexer, I decided that I didnt want to use additional libaries, so I implemented a solution that uses the regexp Go libary to tokenize a string, this given the following rules :

```
	rules["float"] = "([0-9]*[.])+[0-9]+"
	rules["int"] = "[0-9]+"
	rules["var"] = "[a-zA-Z_]+"
	rules["ass"] = "\\="
	rules["add"] = "\\+"
	rules["sub"] = "\\-"
	rules["okey"] = "\\{"
	rules["ckey"] = "\\}"
	rules["opar"] = "\\("

```

I'm currentl working in the defintion of all the rules, alog the way I build de parser

## Parser

To implement the parser, I decided I will use the Recursive Descent Method, with a Context Free Grammar that will prove if all the lines in the source code are valid for the code generator

## Installation

Clone this repository on your computer

```
git clone https://github.com/Bastian1110/Miau65i
```
