CS4513: Project 1 Goat File System
==================================

Note, this document includes a number of design questions that can help your implementation. We highly recommend that you answer each design question **before** attempting the corresponding implementation.
These questions will help you design and plan your implementation and guide you towards the resources you need.
Finally, if you are unsure how to start the project, we recommend you visit office hours for some guidance on these questions before attempting to implement this project.


Team members
-----------------

1. Joshua Malcarne (jrmalcarne@wpi.edu)
2. Miles Gregg (mgregg@wpi.edu)

Design Questions
------------------

1. When implementing the `debug()` function, you will need to load the file system via the emulated disk and retrieve the information for superblock and inodes.
1.1 How will you read the superblock?
1.2 How will you traverse all the inodes?
1.3 How will you determine all the information related to an inode?
1.4 How will you determine all the blocks related to an inode?

Brief response please!

Responses:

    1.1: The reading of the superblock will be done with the wread() built in function 
         from the Disk.c file. Reading from block zero will get the superblock data out
         of the disk. 
    
    1.2: To traverse all of the inodes we plan to use a for loop to iterate through all of
         inode blocks based on the value (number of inode blocks) read from the superblock value from part 1.1.

    1.3: By checking if each of the 128 inodes in each inode block's inode array are valid, and printing the 
         inode's information if it is valid. This includes valid, size, directs, and indirects. We do this for
         each inode block, until we have checked all inodes in a disk.

    1.4: To determine all of the blocks related to an specific inode will need to be done by 
         first checking through the 5 direct blocks first then reading from the indirect (using wread()) if it exists 
         and getting all of it's pointers to other data blocks from the pointer array. 

---

2. When implementing the `format()` function, you will need to write the superblock and clear the remaining blocks in the file system.

2.1 What should happen if the the file system is already mounted?
2.2 What information must be written into the superblock?
2.3 How would you clear all the remaining blocks?

Brief response please!

Responses: 

    2.1: To determine if the file system is already mounted we will do a quick check on the number
         of mounts that the disk has done. So if the value of Mounts (in disk) is not equal to zero 
         then this means that the disk can no longer be formatted (since it has been mounted before) and 
         the file system will throw an appropriate error (the formatting will not be successful).

    2.2: We plan on making a superblock instance and allocating the memory needed by setting each of the 
         following variables to the superblock with appropriate values: magic number, blocks, inode blocks, and inodes.

    2.3: To remove all of the remaining blocks on the disk we plan on making an empty block instance and then writing it
         to the rest of the blocks in the disk. This is dependent on the number of blocks on the disk, as is given by the superblock.

---

3. When implementing the `mount()` function, you will need to prepare a filesystem for use by reading the superblock and allocating the free block bitmap.

3.1 What should happen if the file system is already mounted?
3.2 What sanity checks must you perform before building up the free block bitmaps?
3.3 How will you determine which blocks are free?

Responses: 

    3.1: To determine if the file system is already mounted we will do a quick check on the number
         of mounts that the disk has done. So if the value of Mounts (in disk) is not equal to zero 
         then this means that the disk can no longer be mounted (since it has been mounted before) and 
         the file system will throw an appropriate error (the mounting will not be successful).

    3.2: The sanity checks to first check for are whether the magic number is valid, whether there are a correct
         amount of inodes, and whether the blocks amount detailed by the superblock is equivalent to the number of blocks on the disk.

    3.3: To determine which blocks are free when building our two bitmaps (one for the inodes and one for the data blocks), we will
         loop through all inodes contained on the disk. To start, all data blocks will be set as 0 (free). As we loop through the inodes
         in each inode block, if they are valid we set the corresponding index of the inode bitmap to 1 (not free), otherwise to 0 (free).
         If the inode is valid, we then look at its directs, indirect block, and indirects (after reading in the indirect block with wread())
         to determine which blocks are taken by the inode, and set the corresponding indices of the data block bitmap to 1 (not free) where appropriate.
         Once the bitmaps have been properly built, we can use them in future functions to determine which inodes and which data blocks are free or in-use.

Brief response please!

---

4. To implement `create()`, you will need to locate a free inode and save a new inode into the inode table.

4.1 How will you locate a free inode?
4.2 What information would you see in a new inode?
4.3 How will you record this new inode?

Brief response please!

Responses: 

    4.1: To locate the soonest free inode we will utilize the inode bitmap to iterate through and find
         the first free inode (marked by a 0). If all inodes are in use then we will return -1.

    4.2: In a new inode, you would see an inode.Valid value equal to 1, a size of 0, an empty directs array
         (a directs array with no block assignments yet), and an empty indirect block (no assignment for an indirect block).

    4.3: This new inode will be recording by reading in the inode block, setting the appropriate index in the inode block's
         inode array to the new inode instance created, and then writing the inode block back to the file system (since reads 
         and writes can only be done in 4096 byte chunks). The inode bitmap will also be updated by assigning a 1 to the approprate index.

---

5. To implement `remove()`, you will need to locate the inode and then free its associated blocks.

5.1 How will you determine if the specified inode is valid?
5.2 How will you free the direct blocks?
5.3 How will you free the indirect blocks?
5.4 How will you update the inode table?

Brief response please!

Responses: 

    5.1: To determine if the specificed inode is valid we will check the bitmap to see if the input inode
         number is valid (1, in-use) or not (0, free). Valid meaning that there is something in that position to actually remove.

    5.2: We will free the direct blocks by wiping the block data by creating a new block instance and writing it to each direct,
         before then wiping the inode by writing a new, unvalid inode instance over it, wiping the direct array.

    5.3: We will free the indirect blocks by reading in the indirect block, wiping the data of all indirect data blocks as described
         in 5.2, wiping the indirect block itself the same way, and then wiping the inode as described in 5.2.

    5.4: We will update the inode table (bitmap) by assigning a 0 to the appropriate index after the inode has been properly removed.

---

6. To implement `stat()`, you will need to locate the inode and return its size.

6.1 How will you determine if the specified inode is valid?
6.2 How will you determine the inode's size?

Brief response please!

Responses:

    6.1: To check if the specific inode is valid we will need to go to that specific inumber index
         by doing inumber % 128 after reading in the inode data block to go to the inode, 
         then call the Valid variable inside the Inode typedef struct and seeing if it is set to 1 (valid).

    6.2: Similar to 6.1, by doing the inumber % 128 after reading in the inode data block to get the specific 
         inode then just calling the Size variable inside the Inode typedef struct instead to get the inode's size. 

---

7. To implement `read()`, you will need to locate the inode and copy data from appropriate blocks to the user-specified data buffer.

7.1  How will you determine if the specified inode is valid?
7.2  How will you determine which block to read from?
7.3  How will you handle the offset?
7.4  How will you copy from a block to the data buffer?

Brief response please!

Responses:

    7.1: To check if the specific inode is valid we will need to go to that specific inumber index
         by doing inumber % 128 on the read in inode data block to go to the inode then call the Valid 
         variable inside the Inode typedef struct and seeing if it is set to 1 (valid).

    7.2: To determine which block to read from we will have to determine whether it is an direct or
         indirect block first. To do this we plan on doing offset / 4096 becuase if this is less than
         5 then it means that it is an direct block. Else just read from indirects at ((offset / 4096) - 5) 
         to buffer which indirect pointer is read from in the indirect block.

    7.3: We will calculate offset / 4096 when determining which block to read from, and will index block data
         when copying into the provided data pointer using offset % 4096.

    7.4: We will copy from a block to the data buffer using strncpy. This will copy all data from our indexed data block
         to the data pointer for the specified length passed by the function (so long as we pass it to the strncpy function
         appropriately). We prefer this over memcpy since strncpy will complete and terminate when a null character is encountered.

---

8. To implement `write()`, you will need to locate the inode and copy data the user-specified data buffer to data blocks in the file system.

8.1  How will you determine if the specified inode is valid?
8.2  How will you determine which block to write to?
8.3  How will you handle the offset?
8.4  How will you know if you need a new block?
8.5  How will you manage allocating a new block if you need another one?
8.6  How will you copy to a block from the data buffer?
8.7  How will you update the inode?

Brief response please!

Responses: 

    8.1: To check if the specific inode is valid we will need to go to that specific inumber index
         by doing inumber % 128 on the read in inode data block to go to the inode then call the Valid 
         variable inside the Inode typedef struct and seeing if it is set to 1 (valid).

    8.2: To determine which block to write to we will have to determine whether it is an direct or
         indirect block first. To do this we plan on doing offset / 4096 becuase if this is less than
         5 then it means that it is an direct block. Else just write to indirects at ((offset / 4096) - 5) 
         to buffer which indirect pointer is written to in the indirect block.

    8.3: We will calculate offset / 4096 when determining which block to write into, and will index block data
         when copying from the provided data pointer using offset % 4096. Offset will be updated appropriately
         within our used while loop and will be tracked from its original value using an int prev_offset variable.

    8.4: We will know if we need a new block based on the way we are determining which block to write to: as described in
         8.3, to do this we plan on doing offset / 4096 becuase if this is less than 5 then it means that it is an direct block. 
         Else just write to indirects at ((offset / 4096) - 5) to buffer which indirect pointer is written to in the indirect block. 
         Our while loop will handle copying in data in appropriate increments within these conditions. If we need a new block, one will be allocated.

    8.5: We will use the data block bitmap to determine which data blocks are free, and then assign the data block to the appropriate direct/indirect
         pointer before writing the updated inode/indirect block using wwrite(). This is the case for assigning an indirect block as well (for holding
         indirect pointers). Then, we will update the data block bitmap to 1 (in-use) at the appropriate index and use our newly allocated data block as needed.

    8.6: We will do this using strncpy within our while loop, which will assign appropriately sized chunks of data to be copied at a time. This will copy
         the selected chunk of the data pointer (indexed) to our selected data block (indexed) for the specified length (determined within our loop and passed
         to the strncpy function). We prefer strncpy over memcpy since strncpy will complete and terminate when a null character is encountered. Then, in the
         final data block written to, if the entire block is not filled, we will append a null character to the very end of what is written as a fail-case
         for future strncpy calls in copyout and copyin.

    8.7: The inode will be updated as described in 8.4. After a new block is allocated, we will update directs and indirects as appropriate.
         Then, we will update the inode Size before reading in the inode block, copying the inode to the appropriate index in the inode block's
         inode array, and writing our updated inode block back to the disk.

---


Errata
------

Describe any known errors, bugs, or deviations from the requirements.

- None encountered

---

(Optional) Additional Test Cases
--------------------------------

Describe any new test cases that you developed.

- None developed

