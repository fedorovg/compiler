# Compiler

This is a compiler for a subset of pascal language written in go.

# Dependencies 

I use [llir/llvm](https://github.com/llir/llvm) to emit LLVM ir.

You need to have go 1.16 or later on your system. Go should download the llvm go lib during the build process on its own. 

To actually create an executable from ir via `mila` script you need to have LLVM toolchain installed.

# Examples

## Simple
	
  ```pascal
  program factorial;

function facti(n : integer) : integer;
begin
    facti := 1;
    while n > 1 do
    begin
        facti := facti * n;
        dec(n);
    end
end;    

function factr(n : integer) : integer;
begin
    if n = 1 then 
        factr := 1
    else
        factr := n * factr(n-1);
end;    

begin
    writeln(facti(5));
    writeln(factr(5));
end.

  ```

<details>
  <summary>Emitted IR</summary>
	
  ```assembly
source_filename = "factorial"

declare i32 @writeln(i32 %x)

declare i32 @write(i8* %x)

declare i32 @readln(i32* %x)

define i32 @inc(i32* %x) {
entry:
	%0 = load i32, i32* %x
	%1 = add i32 %0, 1
	store i32 %1, i32* %x
	ret i32 0
}

define i32 @dec(i32* %x) {
entry:
	%0 = load i32, i32* %x
	%1 = sub i32 %0, 1
	store i32 %1, i32* %x
	ret i32 0
}

define i32 @facti(i32 %n) {
entry:
	%0 = alloca i32
	%1 = alloca i32
	store i32 %n, i32* %1
	store i32 1, i32* %0
	%2 = load i32, i32* %1
	%3 = icmp sgt i32 %2, 1
	br i1 %3, label %4, label %11

4:
	%5 = load i32, i32* %0
	%6 = load i32, i32* %1
	%7 = mul i32 %5, %6
	store i32 %7, i32* %0
	%8 = call i32 @dec(i32* %1)
	%9 = load i32, i32* %1
	%10 = icmp sgt i32 %9, 1
	br i1 %10, label %4, label %11

11:
	%12 = load i32, i32* %0
	ret i32 %12
}

define i32 @factr(i32 %n) {
entry:
	%0 = alloca i32
	%1 = alloca i32
	store i32 %n, i32* %1
	%2 = load i32, i32* %1
	%3 = icmp eq i32 %2, 1
	br i1 %3, label %4, label %5

4:
	store i32 1, i32* %0
	br label %11

5:
	%6 = load i32, i32* %1
	%7 = load i32, i32* %1
	%8 = sub i32 %7, 1
	%9 = call i32 @factr(i32 %8)
	%10 = mul i32 %6, %9
	store i32 %10, i32* %0
	br label %11

11:
	%12 = load i32, i32* %0
	ret i32 %12
}

define i32 @main() {
entry:
	%0 = alloca i32
	%1 = call i32 @facti(i32 5)
	%2 = call i32 @writeln(i32 %1)
	%3 = call i32 @factr(i32 5)
	%4 = call i32 @writeln(i32 %3)
	ret i32 0
}


  ```
</details>

## Complex
<details>
  <summary>Input source</summary>
	
  ```pascal
  program gcd;

function gcdi(a: integer; b: integer): integer;
var tmp: integer;
begin
    while b <> 0 do
    begin
        tmp := b;
        b := a mod b;
        a := tmp;
    end;
    gcdi := a;
end;

function gcdr(a: integer; b: integer): integer;
var tmp: integer;
begin
    tmp := a mod b;
    if tmp = 0 then
    begin
        gcdr := b;
        exit;
    end;
    gcdr := gcdr(b, tmp);
end;

function gcdr_guessing_inner(a: integer; b: integer; c: integer): integer;
begin
    if ((a mod c) = 0) and ((b mod c) = 0) then
    begin
        gcdr_guessing_inner := c;
        exit;
    end;
    gcdr_guessing_inner := gcdr_guessing_inner(a, b, c - 1);
end;

function gcdr_guessing(a: integer; b: integer): integer;
begin
    gcdr_guessing := gcdr_guessing_inner(a, b, b);
end;

begin
    writeln(gcdi(27*2, 27*3));
    writeln(gcdr(27*2, 27*3));
    writeln(gcdr_guessing(27*2, 27*3));
end.

  ```
</details>

<details>
  <summary>Emitted IR</summary>
	
  ```assembly
  source_filename = "gcd"

declare i32 @writeln(i32 %x)

declare i32 @write(i8* %x)

declare i32 @readln(i32* %x)

define i32 @inc(i32* %x) {
entry:
	%0 = load i32, i32* %x
	%1 = add i32 %0, 1
	store i32 %1, i32* %x
	ret i32 0
}

define i32 @dec(i32* %x) {
entry:
	%0 = load i32, i32* %x
	%1 = sub i32 %0, 1
	store i32 %1, i32* %x
	ret i32 0
}

define i32 @gcdi(i32 %a, i32 %b) {
entry:
	%0 = alloca i32
	%1 = alloca i32
	store i32 %a, i32* %1
	%2 = alloca i32
	store i32 %b, i32* %2
	%3 = alloca i32
	%4 = load i32, i32* %2
	%5 = icmp ne i32 %4, 0
	br i1 %5, label %6, label %14

6:
	%7 = load i32, i32* %2
	store i32 %7, i32* %3
	%8 = load i32, i32* %1
	%9 = load i32, i32* %2
	%10 = srem i32 %8, %9
	store i32 %10, i32* %2
	%11 = load i32, i32* %3
	store i32 %11, i32* %1
	%12 = load i32, i32* %2
	%13 = icmp ne i32 %12, 0
	br i1 %13, label %6, label %14

14:
	%15 = load i32, i32* %1
	store i32 %15, i32* %0
	%16 = load i32, i32* %0
	ret i32 %16
}

define i32 @gcdr(i32 %a, i32 %b) {
entry:
	%0 = alloca i32
	%1 = alloca i32
	store i32 %a, i32* %1
	%2 = alloca i32
	store i32 %b, i32* %2
	%3 = alloca i32
	%4 = load i32, i32* %1
	%5 = load i32, i32* %2
	%6 = srem i32 %4, %5
	store i32 %6, i32* %3
	%7 = load i32, i32* %3
	%8 = icmp eq i32 %7, 0
	br i1 %8, label %9, label %12

9:
	%10 = load i32, i32* %2
	store i32 %10, i32* %0
	%11 = load i32, i32* %0
	ret i32 %11

12:
	%13 = load i32, i32* %2
	%14 = load i32, i32* %3
	%15 = call i32 @gcdr(i32 %13, i32 %14)
	store i32 %15, i32* %0
	%16 = load i32, i32* %0
	ret i32 %16
}

define i32 @gcdr_guessing_inner(i32 %a, i32 %b, i32 %c) {
entry:
	%0 = alloca i32
	%1 = alloca i32
	store i32 %a, i32* %1
	%2 = alloca i32
	store i32 %b, i32* %2
	%3 = alloca i32
	store i32 %c, i32* %3
	%4 = load i32, i32* %1
	%5 = load i32, i32* %3
	%6 = srem i32 %4, %5
	%7 = icmp eq i32 %6, 0
	%8 = load i32, i32* %2
	%9 = load i32, i32* %3
	%10 = srem i32 %8, %9
	%11 = icmp eq i32 %10, 0
	%12 = and i1 %7, %11
	br i1 %12, label %13, label %16

13:
	%14 = load i32, i32* %3
	store i32 %14, i32* %0
	%15 = load i32, i32* %0
	ret i32 %15

16:
	%17 = load i32, i32* %1
	%18 = load i32, i32* %2
	%19 = load i32, i32* %3
	%20 = sub i32 %19, 1
	%21 = call i32 @gcdr_guessing_inner(i32 %17, i32 %18, i32 %20)
	store i32 %21, i32* %0
	%22 = load i32, i32* %0
	ret i32 %22
}

define i32 @gcdr_guessing(i32 %a, i32 %b) {
entry:
	%0 = alloca i32
	%1 = alloca i32
	store i32 %a, i32* %1
	%2 = alloca i32
	store i32 %b, i32* %2
	%3 = load i32, i32* %1
	%4 = load i32, i32* %2
	%5 = load i32, i32* %2
	%6 = call i32 @gcdr_guessing_inner(i32 %3, i32 %4, i32 %5)
	store i32 %6, i32* %0
	%7 = load i32, i32* %0
	ret i32 %7
}

define i32 @main() {
entry:
	%0 = alloca i32
	%1 = mul i32 27, 2
	%2 = mul i32 27, 3
	%3 = call i32 @gcdi(i32 %1, i32 %2)
	%4 = call i32 @writeln(i32 %3)
	%5 = mul i32 27, 2
	%6 = mul i32 27, 3
	%7 = call i32 @gcdr(i32 %5, i32 %6)
	%8 = call i32 @writeln(i32 %7)
	%9 = mul i32 27, 2
	%10 = mul i32 27, 3
	%11 = call i32 @gcdr_guessing(i32 %9, i32 %10)
	%12 = call i32 @writeln(i32 %11)
	ret i32 0
}


  ```
</details>


# Build

```bash
cd gila
go build -o ../build/gila ./main
```
or just run `make` in the root folder.

# Run

```bash
./mila <file>
```







